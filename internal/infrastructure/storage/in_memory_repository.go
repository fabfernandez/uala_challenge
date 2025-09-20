package storage

import (
	"context"
	"sync"

	"uala-challenge/internal/domain"
)

// InMemoryRepository implements all domain repositories using in-memory storage
type InMemoryRepository struct {
	users   map[string]*domain.User
	tweets  map[string]*domain.Tweet
	follows map[string][]string // followerID -> []followeeID
	mutex   sync.RWMutex
}

// NewInMemoryRepository creates a new in-memory repository
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		users:   make(map[string]*domain.User),
		tweets:  make(map[string]*domain.Tweet),
		follows: make(map[string][]string),
	}
}

// User Repository Implementation

func (r *InMemoryRepository) CreateUser(ctx context.Context, user *domain.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	r.users[user.ID] = user
	return nil
}

func (r *InMemoryRepository) GetUser(ctx context.Context, id string) (*domain.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	user, exists := r.users[id]
	if !exists {
		return nil, nil
	}
	
	return user, nil
}

// Tweet Repository Implementation

func (r *InMemoryRepository) CreateTweet(ctx context.Context, tweet *domain.Tweet) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	r.tweets[tweet.ID] = tweet
	return nil
}

func (r *InMemoryRepository) GetTweetsByUserID(ctx context.Context, userID string) ([]*domain.Tweet, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	var userTweets []*domain.Tweet
	for _, tweet := range r.tweets {
		if tweet.UserID == userID {
			userTweets = append(userTweets, tweet)
		}
	}
	
	return userTweets, nil
}

func (r *InMemoryRepository) GetTweetsByUserIDs(ctx context.Context, userIDs []string) ([]*domain.Tweet, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	userIDSet := make(map[string]bool)
	for _, id := range userIDs {
		userIDSet[id] = true
	}
	
	var tweets []*domain.Tweet
	for _, tweet := range r.tweets {
		if userIDSet[tweet.UserID] {
			tweets = append(tweets, tweet)
		}
	}
	
	return tweets, nil
}

// Follow Repository Implementation

func (r *InMemoryRepository) FollowUser(ctx context.Context, followerID, followeeID string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	// Check if already following
	for _, existingFollowee := range r.follows[followerID] {
		if existingFollowee == followeeID {
			return nil // Already following
		}
	}
	
	r.follows[followerID] = append(r.follows[followerID], followeeID)
	return nil
}

func (r *InMemoryRepository) UnfollowUser(ctx context.Context, followerID, followeeID string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	followees := r.follows[followerID]
	for i, followee := range followees {
		if followee == followeeID {
			// Remove the followee
			r.follows[followerID] = append(followees[:i], followees[i+1:]...)
			break
		}
	}
	
	return nil
}

func (r *InMemoryRepository) GetFollowees(ctx context.Context, followerID string) ([]string, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	followees := r.follows[followerID]
	if followees == nil {
		return []string{}, nil
	}
	
	return followees, nil
}
