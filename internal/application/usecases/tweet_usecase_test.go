package usecases

import (
	"context"
	"testing"

	"uala-challenge/internal/domain"
)

// Mock repositories for testing
type mockUserRepository struct {
	users map[string]*domain.User
}

func (m *mockUserRepository) Create(ctx context.Context, user *domain.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return m.users[id], nil
}

type mockTweetRepository struct {
	tweets []*domain.Tweet
}

func (m *mockTweetRepository) Create(ctx context.Context, tweet *domain.Tweet) error {
	m.tweets = append(m.tweets, tweet)
	return nil
}

func (m *mockTweetRepository) GetByUserID(ctx context.Context, userID string) ([]*domain.Tweet, error) {
	var userTweets []*domain.Tweet
	for _, tweet := range m.tweets {
		if tweet.UserID == userID {
			userTweets = append(userTweets, tweet)
		}
	}
	return userTweets, nil
}

func (m *mockTweetRepository) GetByUserIDs(ctx context.Context, userIDs []string) ([]*domain.Tweet, error) {
	// Not used in tweet use case tests
	return nil, nil
}

func TestTweetUseCase_CreateTweet(t *testing.T) {
	ctx := context.Background()
	
	tests := []struct {
		name        string
		req         CreateTweetRequest
		expectError bool
	}{
		{
			name: "valid tweet",
			req: CreateTweetRequest{
				UserID:  "user123",
				Content: "Hello, world!",
			},
			expectError: false,
		},
		{
			name: "empty content",
			req: CreateTweetRequest{
				UserID:  "user123",
				Content: "",
			},
			expectError: true,
		},
		{
			name: "content too long",
			req: CreateTweetRequest{
				UserID:  "user123",
				Content: string(make([]byte, 281)),
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := &mockUserRepository{users: make(map[string]*domain.User)}
			tweetRepo := &mockTweetRepository{tweets: []*domain.Tweet{}}
			
			useCase := NewTweetUseCase(tweetRepo, userRepo)
			
			tweet, err := useCase.CreateTweet(ctx, tt.req)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if tweet.Content != tt.req.Content {
					t.Errorf("Expected content %s, got %s", tt.req.Content, tweet.Content)
				}
				if tweet.UserID != tt.req.UserID {
					t.Errorf("Expected user ID %s, got %s", tt.req.UserID, tweet.UserID)
				}
			}
		})
	}
}

func TestTweetUseCase_GetUserTweets(t *testing.T) {
	ctx := context.Background()
	userID := "user123"
	
	// Create test tweets
	tweet1 := &domain.Tweet{ID: "1", UserID: userID, Content: "First tweet"}
	tweet2 := &domain.Tweet{ID: "2", UserID: userID, Content: "Second tweet"}
	
	userRepo := &mockUserRepository{users: make(map[string]*domain.User)}
	tweetRepo := &mockTweetRepository{tweets: []*domain.Tweet{tweet1, tweet2}}
	
	useCase := NewTweetUseCase(tweetRepo, userRepo)
	
	tweets, err := useCase.GetUserTweets(ctx, userID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if len(tweets) != 2 {
		t.Errorf("Expected 2 tweets, got %d", len(tweets))
	}
}
