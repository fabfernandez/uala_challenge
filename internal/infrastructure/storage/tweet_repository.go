package storage

import (
	"context"

	"uala-challenge/internal/domain"
)

// TweetRepository implements domain.TweetRepository
type TweetRepository struct {
	storage *InMemoryRepository
}

// NewTweetRepository creates a new tweet repository
func NewTweetRepository(storage *InMemoryRepository) *TweetRepository {
	return &TweetRepository{
		storage: storage,
	}
}

func (r *TweetRepository) Create(ctx context.Context, tweet *domain.Tweet) error {
	return r.storage.CreateTweet(ctx, tweet)
}

func (r *TweetRepository) GetByUserID(ctx context.Context, userID string) ([]*domain.Tweet, error) {
	return r.storage.GetTweetsByUserID(ctx, userID)
}

func (r *TweetRepository) GetByUserIDs(ctx context.Context, userIDs []string) ([]*domain.Tweet, error) {
	return r.storage.GetTweetsByUserIDs(ctx, userIDs)
}
