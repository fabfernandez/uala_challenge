package application

import (
	"context"

	"uala-challenge/internal/application/services"
	"uala-challenge/internal/domain"
)

// TweetServiceInterface defines the interface for tweet services
type TweetServiceInterface interface {
	CreateTweet(ctx context.Context, req services.CreateTweetRequest) (*domain.Tweet, error)
	GetUserTweets(ctx context.Context, userID string) ([]*domain.Tweet, error)
}

// FollowServiceInterface defines the interface for follow services
type FollowServiceInterface interface {
	FollowUser(ctx context.Context, req services.FollowUserRequest) error
	UnfollowUser(ctx context.Context, req services.FollowUserRequest) error
	GetTimeline(ctx context.Context, userID string) ([]*domain.Tweet, error)
}
