package domain

import "context"

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
}

// TweetRepository defines the interface for tweet data operations
type TweetRepository interface {
	Create(ctx context.Context, tweet *Tweet) error
	GetByUserID(ctx context.Context, userID string) ([]*Tweet, error)
	GetByUserIDs(ctx context.Context, userIDs []string) ([]*Tweet, error)
}

// FollowRepository defines the interface for follow relationship operations
type FollowRepository interface {
	Follow(ctx context.Context, followerID, followeeID string) error
	Unfollow(ctx context.Context, followerID, followeeID string) error
	GetFollowees(ctx context.Context, followerID string) ([]string, error)
}
