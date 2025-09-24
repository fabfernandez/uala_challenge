package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"uala-challenge/internal/application/services"
	"uala-challenge/internal/infrastructure/storage"
)

// TestCompleteWorkflow tests the complete microblogging workflow from the demo
func TestCompleteWorkflow(t *testing.T) {
	// Initialize real dependencies (not mocks)
	inMemoryStorage := storage.NewInMemoryRepository()
	userRepo := storage.NewUserRepository(inMemoryStorage)
	tweetRepo := storage.NewTweetRepository(inMemoryStorage)
	followRepo := storage.NewFollowRepository(inMemoryStorage)

	tweetService := services.NewTweetService(tweetRepo, userRepo)
	followService := services.NewFollowService(followRepo, tweetRepo)

	handler := NewHandler(tweetService, followService)
	router := NewRouter(handler)
	httpRouter := router.SetupRoutes()

	// Test 1: Create tweets from different users
	t.Run("Create tweets from multiple users", func(t *testing.T) {
		// Alice's tweet
		req := createTweetRequest("alice123", "Hello, world! This is my first tweet on the microblogging platform. #excited")
		w := httptest.NewRecorder()
		httpRouter.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
		}

		// Bob's tweet
		req = createTweetRequest("bob456", "Hey everyone! Just joined the microblogging platform. Excited to connect with you all! 🚀")
		w = httptest.NewRecorder()
		httpRouter.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
		}

		// Charlie's tweet
		req = createTweetRequest("charlie789", "Just finished reading an amazing book about software architecture. Clean Architecture principles are game-changing! 📚")
		w = httptest.NewRecorder()
		httpRouter.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
		}

		// Dave's tweet (not followed by Alice)
		req = createTweetRequest("dave999", "This tweet should NOT appear in Alice's timeline since she doesn't follow me.")
		w = httptest.NewRecorder()
		httpRouter.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
		}
	})

	// Test 2: Follow relationships
	t.Run("Follow relationships", func(t *testing.T) {
		// Alice follows Bob
		req := createFollowRequest("alice123", "bob456")
		w := httptest.NewRecorder()
		httpRouter.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		// Alice follows Charlie
		req = createFollowRequest("alice123", "charlie789")
		w = httptest.NewRecorder()
		httpRouter.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		// Test self-follow prevention
		req = createFollowRequest("alice123", "alice123")
		w = httptest.NewRecorder()
		httpRouter.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d for self-follow, got %d", http.StatusBadRequest, w.Code)
		}
	})

	// Test 3: Timeline functionality
	t.Run("Timeline shows only followed users", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/timeline", nil)
		req.Header.Set("X-User-ID", "alice123")
		w := httptest.NewRecorder()
		httpRouter.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Failed to unmarshal response: %v", err)
		}

		tweets := response["tweets"].([]interface{})
		count := int(response["count"].(float64))

		// Should have tweets from Bob and Charlie (2 users Alice follows)
		// Alice's own tweet should not be in timeline (she doesn't follow herself)
		// Dave's tweet should not be in timeline (Alice doesn't follow Dave)
		if count < 2 {
			t.Errorf("Expected at least 2 tweets in Alice's timeline, got %d", count)
		}

		// Verify only followed users' tweets appear
		userIDs := make(map[string]bool)
		for _, tweetInterface := range tweets {
			tweet := tweetInterface.(map[string]interface{})
			userID := tweet["user_id"].(string)
			userIDs[userID] = true
		}

		if !userIDs["bob456"] || !userIDs["charlie789"] {
			t.Error("Expected tweets from Bob and Charlie in Alice's timeline")
		}

		if userIDs["dave999"] {
			t.Error("Expected no tweets from Dave (not followed) in Alice's timeline")
		}
	})

	// Test 4: Unfollow functionality
	t.Run("Unfollow and timeline update", func(t *testing.T) {
		// Alice unfollows Charlie
		req := createUnfollowRequest("alice123", "charlie789")
		w := httptest.NewRecorder()
		httpRouter.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		// Check Alice's timeline after unfollowing Charlie
		req = httptest.NewRequest("GET", "/api/v1/timeline", nil)
		req.Header.Set("X-User-ID", "alice123")
		w = httptest.NewRecorder()
		httpRouter.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Failed to unmarshal response: %v", err)
		}

		tweets := response["tweets"].([]interface{})

		// Should only have tweets from Bob now (Charlie was unfollowed)
		userIDs := make(map[string]bool)
		for _, tweetInterface := range tweets {
			tweet := tweetInterface.(map[string]interface{})
			userID := tweet["user_id"].(string)
			userIDs[userID] = true
		}

		if !userIDs["bob456"] {
			t.Error("Expected tweets from Bob in Alice's timeline after unfollowing Charlie")
		}

		if userIDs["charlie789"] {
			t.Error("Expected no tweets from Charlie in Alice's timeline after unfollowing")
		}
	})

	// Test 5: Get specific user's tweets
	t.Run("Get user's tweets", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/users/tweets?user_id=bob456", nil)
		req.Header.Set("X-User-ID", "alice123")
		w := httptest.NewRecorder()
		httpRouter.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Failed to unmarshal response: %v", err)
		}

		count := int(response["count"].(float64))
		if count < 1 {
			t.Errorf("Expected at least 1 tweet from Bob, got %d", count)
		}
	})

	// Test 6: Empty timeline for user with no follows
	t.Run("Empty timeline for user with no follows", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/timeline", nil)
		req.Header.Set("X-User-ID", "bob456") // Bob follows no one
		w := httptest.NewRecorder()
		httpRouter.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Failed to unmarshal response: %v", err)
		}

		count := int(response["count"].(float64))
		if count != 0 {
			t.Errorf("Expected 0 tweets in Bob's timeline (no follows), got %d", count)
		}
	})
}

// TestCharacterLimitEnforcement tests the 280 character limit from the demo
func TestCharacterLimitEnforcement(t *testing.T) {
	inMemoryStorage := storage.NewInMemoryRepository()
	userRepo := storage.NewUserRepository(inMemoryStorage)
	tweetRepo := storage.NewTweetRepository(inMemoryStorage)
	followRepo := storage.NewFollowRepository(inMemoryStorage)

	tweetService := services.NewTweetService(tweetRepo, userRepo)
	followService := services.NewFollowService(followRepo, tweetRepo)

	handler := NewHandler(tweetService, followService)
	router := NewRouter(handler)
	httpRouter := router.SetupRoutes()

	t.Run("Valid tweet within limit", func(t *testing.T) {
		req := createTweetRequest("user123", "This is a valid tweet within the 280 character limit.")
		w := httptest.NewRecorder()
		httpRouter.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
		}
	})

	t.Run("Tweet exceeding character limit", func(t *testing.T) {
		longTweet := "This tweet is way too long and should be rejected by the system because it exceeds the 280 character limit that was specified in the requirements. Let me keep typing to make sure this definitely goes over the limit. This is a very long tweet that should definitely exceed 280 characters and get rejected by our validation system. I am still typing to make sure this is long enough to trigger the character limit validation. This should definitely be over 280 characters now and should be rejected."

		req := createTweetRequest("user123", longTweet)
		w := httptest.NewRecorder()
		httpRouter.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d for tweet exceeding limit, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("Empty tweet content", func(t *testing.T) {
		req := createTweetRequest("user123", "")
		w := httptest.NewRecorder()
		httpRouter.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d for empty tweet, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

// TestErrorHandling tests error scenarios from the demo
func TestErrorHandling(t *testing.T) {
	inMemoryStorage := storage.NewInMemoryRepository()
	userRepo := storage.NewUserRepository(inMemoryStorage)
	tweetRepo := storage.NewTweetRepository(inMemoryStorage)
	followRepo := storage.NewFollowRepository(inMemoryStorage)

	tweetService := services.NewTweetService(tweetRepo, userRepo)
	followService := services.NewFollowService(followRepo, tweetRepo)

	handler := NewHandler(tweetService, followService)
	router := NewRouter(handler)
	httpRouter := router.SetupRoutes()

	t.Run("Missing user ID header", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/tweets", bytes.NewBufferString(`{"content": "This should fail"}`))
		req.Header.Set("Content-Type", "application/json")
		// No X-User-ID header
		w := httptest.NewRecorder()
		httpRouter.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d for missing user ID, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/tweets", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", "user123")
		w := httptest.NewRecorder()
		httpRouter.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d for invalid JSON, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

// Helper functions
func createTweetRequest(userID, content string) *http.Request {
	reqBody := CreateTweetRequest{Content: content}
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/v1/tweets", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", userID)
	return req
}

func createFollowRequest(followerID, followeeID string) *http.Request {
	reqBody := FollowUserRequest{FolloweeID: followeeID}
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/v1/follow", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", followerID)
	return req
}

func createUnfollowRequest(followerID, followeeID string) *http.Request {
	reqBody := FollowUserRequest{FolloweeID: followeeID}
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/v1/unfollow", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", followerID)
	return req
}
