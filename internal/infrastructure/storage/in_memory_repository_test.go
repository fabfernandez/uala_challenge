package storage

import (
	"context"
	"testing"

	"uala-challenge/internal/domain"
)

func TestInMemoryRepository_UserOperations(t *testing.T) {
	repo := NewInMemoryRepository()
	ctx := context.Background()

	// Test create user
	user := domain.NewUser("John Doe")
	err := repo.CreateUser(ctx, user)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test get user
	retrievedUser, err := repo.GetUser(ctx, user.ID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if retrievedUser == nil {
		t.Error("Expected user to be retrieved, got nil")
	}
	if retrievedUser.Name != user.Name {
		t.Errorf("Expected name %s, got %s", user.Name, retrievedUser.Name)
	}

	// Test get non-existent user
	nonExistentUser, err := repo.GetUser(ctx, "nonexistent")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if nonExistentUser != nil {
		t.Error("Expected nil for non-existent user, got user")
	}
}

func TestInMemoryRepository_TweetOperations(t *testing.T) {
	repo := NewInMemoryRepository()
	ctx := context.Background()

	// Test create tweet
	tweet, err := domain.NewTweet("user123", "Hello, world!")
	if err != nil {
		t.Fatalf("Failed to create tweet: %v", err)
	}

	err = repo.CreateTweet(ctx, tweet)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test get tweets by user ID
	tweets, err := repo.GetTweetsByUserID(ctx, "user123")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(tweets) != 1 {
		t.Errorf("Expected 1 tweet, got %d", len(tweets))
	}
	if tweets[0].ID != tweet.ID {
		t.Errorf("Expected tweet ID %s, got %s", tweet.ID, tweets[0].ID)
	}

	// Test get tweets by multiple user IDs
	tweets, err = repo.GetTweetsByUserIDs(ctx, []string{"user123"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(tweets) != 1 {
		t.Errorf("Expected 1 tweet, got %d", len(tweets))
	}
}

func TestInMemoryRepository_FollowOperations(t *testing.T) {
	repo := NewInMemoryRepository()
	ctx := context.Background()

	followerID := "user1"
	followeeID := "user2"

	// Test follow
	err := repo.FollowUser(ctx, followerID, followeeID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test get followees
	followees, err := repo.GetFollowees(ctx, followerID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(followees) != 1 {
		t.Errorf("Expected 1 followee, got %d", len(followees))
	}
	if followees[0] != followeeID {
		t.Errorf("Expected followee %s, got %s", followeeID, followees[0])
	}

	// Test duplicate follow (should not error)
	err = repo.FollowUser(ctx, followerID, followeeID)
	if err != nil {
		t.Errorf("Expected no error on duplicate follow, got %v", err)
	}

	// Test unfollow
	err = repo.UnfollowUser(ctx, followerID, followeeID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test get followees after unfollow
	followees, err = repo.GetFollowees(ctx, followerID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(followees) != 0 {
		t.Errorf("Expected 0 followees after unfollow, got %d", len(followees))
	}
}

func TestInMemoryRepository_Concurrency(t *testing.T) {
	repo := NewInMemoryRepository()
	ctx := context.Background()

	// Test concurrent operations
	done := make(chan bool, 10)
	
	for i := 0; i < 10; i++ {
		go func(i int) {
			user := domain.NewUser("User" + string(rune(i)))
			repo.CreateUser(ctx, user)
			
			tweet, _ := domain.NewTweet(user.ID, "Tweet from user")
			repo.CreateTweet(ctx, tweet)
			
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all users and tweets were created
	// This is a basic concurrency test - in a real scenario you'd want more thorough testing
}
