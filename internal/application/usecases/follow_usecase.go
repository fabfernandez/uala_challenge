package usecases

import (
	"context"
	"sort"

	"uala-challenge/internal/domain"
)

// FollowUseCase handles follow-related business logic
type FollowUseCase struct {
	followRepo domain.FollowRepository
	tweetRepo  domain.TweetRepository
}

// NewFollowUseCase creates a new follow use case
func NewFollowUseCase(followRepo domain.FollowRepository, tweetRepo domain.TweetRepository) *FollowUseCase {
	return &FollowUseCase{
		followRepo: followRepo,
		tweetRepo:  tweetRepo,
	}
}

// FollowUserRequest represents the request to follow a user
type FollowUserRequest struct {
	FollowerID string `json:"follower_id"`
	FolloweeID string `json:"followee_id"`
}

// FollowUser creates a follow relationship
func (uc *FollowUseCase) FollowUser(ctx context.Context, req FollowUserRequest) error {
	// Validate follow relationship
	err := domain.ValidateFollow(req.FollowerID, req.FolloweeID)
	if err != nil {
		return err
	}

	// Create follow relationship
	return uc.followRepo.Follow(ctx, req.FollowerID, req.FolloweeID)
}

// UnfollowUser removes a follow relationship
func (uc *FollowUseCase) UnfollowUser(ctx context.Context, req FollowUserRequest) error {
	return uc.followRepo.Unfollow(ctx, req.FollowerID, req.FolloweeID)
}

// GetTimeline retrieves tweets from followed users
func (uc *FollowUseCase) GetTimeline(ctx context.Context, userID string) ([]*domain.Tweet, error) {
	// Get list of followed users
	followees, err := uc.followRepo.GetFollowees(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(followees) == 0 {
		return []*domain.Tweet{}, nil
	}

	// Get tweets from followed users
	tweets, err := uc.tweetRepo.GetByUserIDs(ctx, followees)
	if err != nil {
		return nil, err
	}

	// Sort by creation time (newest first)
	sort.Slice(tweets, func(i, j int) bool {
		return tweets[i].CreatedAt.After(tweets[j].CreatedAt)
	})

	return tweets, nil
}
