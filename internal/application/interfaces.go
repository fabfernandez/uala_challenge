package application

import (
	"context"

	"uala-challenge/internal/application/usecases"
	"uala-challenge/internal/domain"
)

// TweetUseCaseInterface defines the interface for tweet use cases
type TweetUseCaseInterface interface {
	CreateTweet(ctx context.Context, req usecases.CreateTweetRequest) (*domain.Tweet, error)
	GetUserTweets(ctx context.Context, userID string) ([]*domain.Tweet, error)
}

// FollowUseCaseInterface defines the interface for follow use cases
type FollowUseCaseInterface interface {
	FollowUser(ctx context.Context, req usecases.FollowUserRequest) error
	UnfollowUser(ctx context.Context, req usecases.FollowUserRequest) error
	GetTimeline(ctx context.Context, userID string) ([]*domain.Tweet, error)
}
