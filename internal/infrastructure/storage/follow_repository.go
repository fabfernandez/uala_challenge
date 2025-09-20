package storage

import (
	"context"
)

// FollowRepository implements domain.FollowRepository
type FollowRepository struct {
	storage *InMemoryRepository
}

// NewFollowRepository creates a new follow repository
func NewFollowRepository(storage *InMemoryRepository) *FollowRepository {
	return &FollowRepository{
		storage: storage,
	}
}

func (r *FollowRepository) Follow(ctx context.Context, followerID, followeeID string) error {
	return r.storage.FollowUser(ctx, followerID, followeeID)
}

func (r *FollowRepository) Unfollow(ctx context.Context, followerID, followeeID string) error {
	return r.storage.UnfollowUser(ctx, followerID, followeeID)
}

func (r *FollowRepository) GetFollowees(ctx context.Context, followerID string) ([]string, error) {
	return r.storage.GetFollowees(ctx, followerID)
}
