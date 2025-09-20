package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"uala-challenge/internal/application/usecases"
	"uala-challenge/internal/domain"
)

// Mock use cases for testing
type mockTweetUseCase struct{}

func (m *mockTweetUseCase) CreateTweet(ctx context.Context, req usecases.CreateTweetRequest) (*domain.Tweet, error) {
	if req.Content == "" {
		return nil, domain.ErrTweetEmpty
	}
	if len(req.Content) > domain.MaxTweetLength {
		return nil, domain.ErrTweetTooLong
	}
	return &domain.Tweet{
		ID:      "tweet123",
		UserID:  req.UserID,
		Content: req.Content,
	}, nil
}

func (m *mockTweetUseCase) GetUserTweets(ctx context.Context, userID string) ([]*domain.Tweet, error) {
	return []*domain.Tweet{
		{ID: "1", UserID: userID, Content: "Test tweet"},
	}, nil
}

type mockFollowUseCase struct{}

func (m *mockFollowUseCase) FollowUser(ctx context.Context, req usecases.FollowUserRequest) error {
	if req.FollowerID == req.FolloweeID {
		return domain.ErrCannotFollowSelf
	}
	return nil
}

func (m *mockFollowUseCase) UnfollowUser(ctx context.Context, req usecases.FollowUserRequest) error {
	return nil
}

func (m *mockFollowUseCase) GetTimeline(ctx context.Context, userID string) ([]*domain.Tweet, error) {
	return []*domain.Tweet{
		{ID: "1", UserID: "other", Content: "Timeline tweet"},
	}, nil
}

func TestHandler_CreateTweetHandler(t *testing.T) {
	handler := NewHandler(&mockTweetUseCase{}, &mockFollowUseCase{})

	tests := []struct {
		name           string
		userID         string
		content        string
		expectedStatus int
	}{
		{
			name:           "valid tweet",
			userID:         "user123",
			content:        "Hello, world!",
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "empty content",
			userID:         "user123",
			content:        "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "missing user ID",
			userID:         "",
			content:        "Hello, world!",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody := CreateTweetRequest{Content: tt.content}
			jsonBody, _ := json.Marshal(reqBody)

			req := httptest.NewRequest("POST", "/api/v1/tweets", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			if tt.userID != "" {
				req.Header.Set("X-User-ID", tt.userID)
			}

			w := httptest.NewRecorder()
			handler.CreateTweetHandler(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestHandler_GetTimelineHandler(t *testing.T) {
	handler := NewHandler(&mockTweetUseCase{}, &mockFollowUseCase{})

	req := httptest.NewRequest("GET", "/api/v1/timeline", nil)
	req.Header.Set("X-User-ID", "user123")

	w := httptest.NewRecorder()
	handler.GetTimelineHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response["count"] != float64(1) {
		t.Errorf("Expected count 1, got %v", response["count"])
	}
}

func TestHandler_FollowUserHandler(t *testing.T) {
	handler := NewHandler(&mockTweetUseCase{}, &mockFollowUseCase{})

	tests := []struct {
		name           string
		userID         string
		followeeID     string
		expectedStatus int
	}{
		{
			name:           "valid follow",
			userID:         "user1",
			followeeID:     "user2",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "self follow",
			userID:         "user1",
			followeeID:     "user1",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "missing user ID",
			userID:         "",
			followeeID:     "user2",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody := FollowUserRequest{FolloweeID: tt.followeeID}
			jsonBody, _ := json.Marshal(reqBody)

			req := httptest.NewRequest("POST", "/api/v1/follow", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			if tt.userID != "" {
				req.Header.Set("X-User-ID", tt.userID)
			}

			w := httptest.NewRecorder()
			handler.FollowUserHandler(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestHandler_GetUserTweetsHandler(t *testing.T) {
	handler := NewHandler(&mockTweetUseCase{}, &mockFollowUseCase{})

	tests := []struct {
		name           string
		targetUserID   string
		expectedStatus int
		expectedCount  int
	}{
		{
			name:           "valid user tweets request",
			targetUserID:   "user123",
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
		{
			name:           "missing user ID parameter",
			targetUserID:   "",
			expectedStatus: http.StatusBadRequest,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/api/v1/users/tweets"
			if tt.targetUserID != "" {
				url += "?user_id=" + tt.targetUserID
			}

			req := httptest.NewRequest("GET", url, nil)
			req.Header.Set("X-User-ID", "requesting_user")

			w := httptest.NewRecorder()
			handler.GetUserTweetsHandler(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}

				count := int(response["count"].(float64))
				if count != tt.expectedCount {
					t.Errorf("Expected %d tweets, got %d", tt.expectedCount, count)
				}
			}
		})
	}
}

func TestHandler_UnfollowUserHandler(t *testing.T) {
	handler := NewHandler(&mockTweetUseCase{}, &mockFollowUseCase{})

	tests := []struct {
		name           string
		userID         string
		followeeID     string
		expectedStatus int
	}{
		{
			name:           "valid unfollow request",
			userID:         "user1",
			followeeID:     "user2",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "missing user ID",
			userID:         "",
			followeeID:     "user2",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "missing followee ID",
			userID:         "user1",
			followeeID:     "",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody := FollowUserRequest{FolloweeID: tt.followeeID}
			jsonBody, _ := json.Marshal(reqBody)

			req := httptest.NewRequest("POST", "/api/v1/unfollow", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			if tt.userID != "" {
				req.Header.Set("X-User-ID", tt.userID)
			}

			w := httptest.NewRecorder()
			handler.UnfollowUserHandler(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}

				if response["message"] != "Successfully unfollowed user" {
					t.Errorf("Expected success message, got %s", response["message"])
				}
			}
		})
	}
}

func TestHandler_HealthCheckHandler(t *testing.T) {
	handler := NewHandler(&mockTweetUseCase{}, &mockFollowUseCase{})

	req := httptest.NewRequest("GET", "/api/v1/health", nil)
	w := httptest.NewRecorder()

	handler.HealthCheckHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got %s", response["status"])
	}
}
