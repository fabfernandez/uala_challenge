package usecases

import (
	"context"
	"sort"

	"uala-challenge/internal/domain"
)

// TweetUseCase handles tweet-related business logic
type TweetUseCase struct {
	tweetRepo domain.TweetRepository
	userRepo  domain.UserRepository
}

// NewTweetUseCase creates a new tweet use case
func NewTweetUseCase(tweetRepo domain.TweetRepository, userRepo domain.UserRepository) *TweetUseCase {
	return &TweetUseCase{
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
func (uc *TweetUseCase) CreateTweet(ctx context.Context, req CreateTweetRequest) (*domain.Tweet, error) {
	// Check if user exists, create if not
	user, err := uc.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	
	if user == nil {
		// Create a default user for this ID
		user = domain.NewUser("User-" + req.UserID)
		err = uc.userRepo.Create(ctx, user)
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
	err = uc.tweetRepo.Create(ctx, tweet)
	if err != nil {
		return nil, err
	}

	return tweet, nil
}

// GetUserTweets retrieves all tweets for a specific user
func (uc *TweetUseCase) GetUserTweets(ctx context.Context, userID string) ([]*domain.Tweet, error) {
	tweets, err := uc.tweetRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Sort by creation time (newest first)
	sort.Slice(tweets, func(i, j int) bool {
		return tweets[i].CreatedAt.After(tweets[j].CreatedAt)
	})

	return tweets, nil
}
