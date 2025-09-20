package domain

import (
	"testing"
)

func TestNewUser(t *testing.T) {
	name := "John Doe"

	user := NewUser(name)

	if user.ID == "" {
		t.Error("User ID should not be empty")
	}

	if user.Name != name {
		t.Errorf("Expected Name %s, got %s", name, user.Name)
	}
}

func TestNewTweet(t *testing.T) {
	userID := "user123"
	
	tests := []struct {
		name        string
		content     string
		expectError bool
		errorType   error
	}{
		{
			name:        "valid tweet",
			content:     "Hello, world!",
			expectError: false,
		},
		{
			name:        "empty content",
			content:     "",
			expectError: true,
			errorType:   ErrTweetEmpty,
		},
		{
			name:        "whitespace only content",
			content:     "   ",
			expectError: true,
			errorType:   ErrTweetEmpty,
		},
		{
			name:        "content exceeding character limit",
			content:     string(make([]byte, MaxTweetLength+1)),
			expectError: true,
			errorType:   ErrTweetTooLong,
		},
		{
			name:        "exactly 280 characters",
			content:     string(make([]byte, MaxTweetLength)),
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tweet, err := NewTweet(userID, tt.content)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
				if err != tt.errorType {
					t.Errorf("Expected error %v, got %v", tt.errorType, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if tweet.UserID != userID {
					t.Errorf("Expected UserID %s, got %s", userID, tweet.UserID)
				}
				if tweet.Content != tt.content {
					t.Errorf("Expected Content %s, got %s", tt.content, tweet.Content)
				}
				if tweet.ID == "" {
					t.Error("Tweet ID should not be empty")
				}
			}
		})
	}
}

func TestValidateFollow(t *testing.T) {
	tests := []struct {
		name        string
		followerID  string
		followeeID  string
		expectError bool
		errorType   error
	}{
		{
			name:        "valid follow",
			followerID:  "user1",
			followeeID:  "user2",
			expectError: false,
		},
		{
			name:        "self follow",
			followerID:  "user1",
			followeeID:  "user1",
			expectError: true,
			errorType:   ErrCannotFollowSelf,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFollow(tt.followerID, tt.followeeID)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
				if err != tt.errorType {
					t.Errorf("Expected error %v, got %v", tt.errorType, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			}
		})
	}
}
