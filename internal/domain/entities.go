package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Domain errors
var (
	ErrTweetTooLong     = errors.New("tweet exceeds character limit")
	ErrTweetEmpty       = errors.New("tweet content cannot be empty")
	ErrUserNotFound     = errors.New("user not found")
	ErrCannotFollowSelf = errors.New("cannot follow yourself")
)

const MaxTweetLength = 280

// User represents a user in the system
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Tweet represents a tweet/post
type Tweet struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// Follow represents a follow relationship between users
type Follow struct {
	FollowerID string `json:"follower_id"`
	FolloweeID string `json:"followee_id"`
}

// NewUser creates a new user with generated ID
func NewUser(name string) *User {
	return &User{
		ID:   uuid.New().String(),
		Name: name,
	}
}

// NewTweet creates a new tweet with validation
func NewTweet(userID, content string) (*Tweet, error) {
	// Validate content
	if len(strings.TrimSpace(content)) == 0 {
		return nil, ErrTweetEmpty
	}
	
	if len(content) > MaxTweetLength {
		return nil, ErrTweetTooLong
	}

	return &Tweet{
		ID:        uuid.New().String(),
		UserID:    userID,
		Content:   content,
		CreatedAt: time.Now(),
	}, nil
}

// ValidateFollow checks if a follow relationship is valid
func ValidateFollow(followerID, followeeID string) error {
	if followerID == followeeID {
		return ErrCannotFollowSelf
	}
	return nil
}
