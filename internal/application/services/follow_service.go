package services

import (
	"context"
	"sort"

	"uala-challenge/internal/domain"
)

// FollowService handles follow-related business logic
type FollowService struct {
	followRepo domain.FollowRepository
	tweetRepo  domain.TweetRepository
}

// NewFollowService creates a new follow service
func NewFollowService(followRepo domain.FollowRepository, tweetRepo domain.TweetRepository) *FollowService {
	return &FollowService{
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
func (s *FollowService) FollowUser(ctx context.Context, req FollowUserRequest) error {
	// Validate follow relationship
	err := domain.ValidateFollow(req.FollowerID, req.FolloweeID)
	if err != nil {
		return err
	}

	// Create follow relationship
	return s.followRepo.Follow(ctx, req.FollowerID, req.FolloweeID)
}

// UnfollowUser removes a follow relationship
func (s *FollowService) UnfollowUser(ctx context.Context, req FollowUserRequest) error {
	return s.followRepo.Unfollow(ctx, req.FollowerID, req.FolloweeID)
}

// GetTimeline retrieves tweets from followed users
func (s *FollowService) GetTimeline(ctx context.Context, userID string) ([]*domain.Tweet, error) {
	// Get list of followed users
	followees, err := s.followRepo.GetFollowees(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(followees) == 0 {
		return []*domain.Tweet{}, nil
	}

	// Get tweets from followed users
	tweets, err := s.tweetRepo.GetByUserIDs(ctx, followees)
	if err != nil {
		return nil, err
	}

	// Sort by creation time (newest first)
	sort.Slice(tweets, func(i, j int) bool {
		return tweets[i].CreatedAt.After(tweets[j].CreatedAt)
	})

	return tweets, nil
}
