package services

import (
	"context"
	"testing"

	"uala-challenge/internal/domain"
)

// Mock repositories for testing
type mockFollowRepository struct {
	follows map[string][]string // followerID -> []followeeID
}

func (m *mockFollowRepository) Follow(ctx context.Context, followerID, followeeID string) error {
	m.follows[followerID] = append(m.follows[followerID], followeeID)
	return nil
}

func (m *mockFollowRepository) Unfollow(ctx context.Context, followerID, followeeID string) error {
	followees := m.follows[followerID]
	for i, followee := range followees {
		if followee == followeeID {
			m.follows[followerID] = append(followees[:i], followees[i+1:]...)
			break
		}
	}
	return nil
}

func (m *mockFollowRepository) GetFollowees(ctx context.Context, followerID string) ([]string, error) {
	return m.follows[followerID], nil
}

type mockTweetRepositoryForFollow struct {
	tweets []*domain.Tweet
}

func (m *mockTweetRepositoryForFollow) Create(ctx context.Context, tweet *domain.Tweet) error {
	m.tweets = append(m.tweets, tweet)
	return nil
}

func (m *mockTweetRepositoryForFollow) GetByUserID(ctx context.Context, userID string) ([]*domain.Tweet, error) {
	var userTweets []*domain.Tweet
	for _, tweet := range m.tweets {
		if tweet.UserID == userID {
			userTweets = append(userTweets, tweet)
		}
	}
	return userTweets, nil
}

func (m *mockTweetRepositoryForFollow) GetByUserIDs(ctx context.Context, userIDs []string) ([]*domain.Tweet, error) {
	userIDSet := make(map[string]bool)
	for _, id := range userIDs {
		userIDSet[id] = true
	}

	var tweets []*domain.Tweet
	for _, tweet := range m.tweets {
		if userIDSet[tweet.UserID] {
			tweets = append(tweets, tweet)
		}
	}
	return tweets, nil
}

func TestFollowService_FollowUser(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		req         FollowUserRequest
		expectError bool
	}{
		{
			name: "valid follow",
			req: FollowUserRequest{
				FollowerID: "user1",
				FolloweeID: "user2",
			},
			expectError: false,
		},
		{
			name: "self follow",
			req: FollowUserRequest{
				FollowerID: "user1",
				FolloweeID: "user1",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			followRepo := &mockFollowRepository{follows: make(map[string][]string)}
			tweetRepo := &mockTweetRepositoryForFollow{tweets: []*domain.Tweet{}}

			service := NewFollowService(followRepo, tweetRepo)

			err := service.FollowUser(ctx, tt.req)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			}
		})
	}
}

func TestFollowService_UnfollowUser(t *testing.T) {
	ctx := context.Background()

	followRepo := &mockFollowRepository{
		follows: map[string][]string{
			"user1": {"user2", "user3"},
		},
	}
	tweetRepo := &mockTweetRepositoryForFollow{tweets: []*domain.Tweet{}}

	service := NewFollowService(followRepo, tweetRepo)

	req := FollowUserRequest{
		FollowerID: "user1",
		FolloweeID: "user2",
	}

	err := service.UnfollowUser(ctx, req)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify unfollow worked
	followees, err := followRepo.GetFollowees(ctx, "user1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(followees) != 1 {
		t.Errorf("Expected 1 followee after unfollow, got %d", len(followees))
	}

	if followees[0] != "user3" {
		t.Errorf("Expected remaining followee to be user3, got %s", followees[0])
	}
}

func TestFollowService_GetTimeline(t *testing.T) {
	ctx := context.Background()

	followRepo := &mockFollowRepository{
		follows: map[string][]string{
			"user1": {"user2", "user3"},
		},
	}

	// Create test tweets
	tweet1 := &domain.Tweet{ID: "1", UserID: "user2", Content: "Tweet from user2"}
	tweet2 := &domain.Tweet{ID: "2", UserID: "user3", Content: "Tweet from user3"}
	tweet3 := &domain.Tweet{ID: "3", UserID: "user4", Content: "Tweet from user4 (not followed)"}

	tweetRepo := &mockTweetRepositoryForFollow{
		tweets: []*domain.Tweet{tweet1, tweet2, tweet3},
	}

	service := NewFollowService(followRepo, tweetRepo)

	tweets, err := service.GetTimeline(ctx, "user1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(tweets) != 2 {
		t.Errorf("Expected 2 tweets in timeline, got %d", len(tweets))
	}

	// Verify only followed users' tweets are included
	userIDs := make(map[string]bool)
	for _, tweet := range tweets {
		userIDs[tweet.UserID] = true
	}

	if !userIDs["user2"] || !userIDs["user3"] {
		t.Error("Expected tweets from user2 and user3 in timeline")
	}

	if userIDs["user4"] {
		t.Error("Expected no tweets from user4 (not followed) in timeline")
	}
}

func TestFollowService_GetTimeline_EmptyFollows(t *testing.T) {
	ctx := context.Background()

	followRepo := &mockFollowRepository{
		follows: map[string][]string{},
	}
	tweetRepo := &mockTweetRepositoryForFollow{tweets: []*domain.Tweet{}}

	service := NewFollowService(followRepo, tweetRepo)

	tweets, err := service.GetTimeline(ctx, "user1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(tweets) != 0 {
		t.Errorf("Expected 0 tweets for user with no follows, got %d", len(tweets))
	}
}
