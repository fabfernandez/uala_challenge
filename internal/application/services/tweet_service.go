package services

import (
	"context"
	"sort"

	"uala-challenge/internal/domain"
)

// TweetService handles tweet-related business logic
type TweetService struct {
	tweetRepo domain.TweetRepository
	userRepo  domain.UserRepository
}

// NewTweetService creates a new tweet service
func NewTweetService(tweetRepo domain.TweetRepository, userRepo domain.UserRepository) *TweetService {
	return &TweetService{
		tweetRepo: tweetRepo,
		userRepo:  userRepo,
	}
}

// CreateTweetRequest represents the request to create a tweet
type CreateTweetRequest struct {
	UserID  string `json:"user_id"`
	Content string `json:"content"`
}

// CreateTweet creates a new tweet
func (s *TweetService) CreateTweet(ctx context.Context, req CreateTweetRequest) (*domain.Tweet, error) {
	// Check if user exists, create if not
	user, err := s.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		// Create a default user for this ID
		user = domain.NewUser("User-" + req.UserID)
		err = s.userRepo.Create(ctx, user)
		if err != nil {
			return nil, err
		}
	}

	// Create tweet with domain validation
	tweet, err := domain.NewTweet(req.UserID, req.Content)
	if err != nil {
		return nil, err
	}

	// Save tweet
	err = s.tweetRepo.Create(ctx, tweet)
	if err != nil {
		return nil, err
	}

	return tweet, nil
}

// GetUserTweets retrieves all tweets for a specific user
func (s *TweetService) GetUserTweets(ctx context.Context, userID string) ([]*domain.Tweet, error) {
	tweets, err := s.tweetRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Sort by creation time (newest first)
	sort.Slice(tweets, func(i, j int) bool {
		return tweets[i].CreatedAt.After(tweets[j].CreatedAt)
	})

	return tweets, nil
}
