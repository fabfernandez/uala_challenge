package http

import (
	"encoding/json"
	"net/http"

	"uala-challenge/internal/application"
	"uala-challenge/internal/application/usecases"
	"uala-challenge/internal/domain"
)

// Handler handles HTTP requests
type Handler struct {
	tweetUseCase  application.TweetUseCaseInterface
	followUseCase application.FollowUseCaseInterface
}

func NewHandler(tweetUseCase application.TweetUseCaseInterface, followUseCase application.FollowUseCaseInterface) *Handler {
	return &Handler{
		tweetUseCase:  tweetUseCase,
		followUseCase: followUseCase,
	}
}

type CreateTweetRequest struct {
	Content string `json:"content"`
}

type FollowUserRequest struct {
	FolloweeID string `json:"followee_id"`
}

func (h *Handler) CreateTweetHandler(w http.ResponseWriter, r *http.Request) {

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "User ID required in X-User-ID header", http.StatusBadRequest)
		return
	}

	var req CreateTweetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	tweet, err := h.tweetUseCase.CreateTweet(r.Context(), usecases.CreateTweetRequest{
		UserID:  userID,
		Content: req.Content,
	})

	if err != nil {
		switch err {
		case domain.ErrTweetEmpty:
			http.Error(w, "Tweet content cannot be empty", http.StatusBadRequest)
		case domain.ErrTweetTooLong:
			http.Error(w, "Tweet content exceeds character limit", http.StatusBadRequest)
		default:
			http.Error(w, "Failed to create tweet", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tweet)
}

func (h *Handler) GetTimelineHandler(w http.ResponseWriter, r *http.Request) {

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "User ID required in X-User-ID header", http.StatusBadRequest)
		return
	}

	tweets, err := h.followUseCase.GetTimeline(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to get timeline", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"tweets": tweets,
		"count":  len(tweets),
	})
}

func (h *Handler) GetUserTweetsHandler(w http.ResponseWriter, r *http.Request) {

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	tweets, err := h.tweetUseCase.GetUserTweets(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to get user tweets", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"tweets": tweets,
		"count":  len(tweets),
	})
}

func (h *Handler) FollowUserHandler(w http.ResponseWriter, r *http.Request) {

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "User ID required in X-User-ID header", http.StatusBadRequest)
		return
	}

	var req FollowUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.FolloweeID == "" {
		http.Error(w, "Followee ID required", http.StatusBadRequest)
		return
	}

	err := h.followUseCase.FollowUser(r.Context(), usecases.FollowUserRequest{
		FollowerID: userID,
		FolloweeID: req.FolloweeID,
	})

	if err != nil {
		switch err {
		case domain.ErrCannotFollowSelf:
			http.Error(w, "Cannot follow yourself", http.StatusBadRequest)
		default:
			http.Error(w, "Failed to follow user", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Successfully followed user",
	})
}

func (h *Handler) UnfollowUserHandler(w http.ResponseWriter, r *http.Request) {

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "User ID required in X-User-ID header", http.StatusBadRequest)
		return
	}

	var req FollowUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.FolloweeID == "" {
		http.Error(w, "Followee ID required", http.StatusBadRequest)
		return
	}

	err := h.followUseCase.UnfollowUser(r.Context(), usecases.FollowUserRequest{
		FollowerID: userID,
		FolloweeID: req.FolloweeID,
	})

	if err != nil {
		http.Error(w, "Failed to unfollow user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Successfully unfollowed user",
	})
}

func (h *Handler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
	})
}
