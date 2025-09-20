package http

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Router sets up HTTP routes
type Router struct {
	handler *Handler
}

// NewRouter creates a new router
func NewRouter(handler *Handler) *Router {
	return &Router{
		handler: handler,
	}
}

// SetupRoutes configures all HTTP routes
func (r *Router) SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// Tweet routes
	api.HandleFunc("/tweets", r.handler.CreateTweetHandler).Methods("POST")
	api.HandleFunc("/timeline", r.handler.GetTimelineHandler).Methods("GET")
	api.HandleFunc("/users/tweets", r.handler.GetUserTweetsHandler).Methods("GET")

	// Follow routes
	api.HandleFunc("/follow", r.handler.FollowUserHandler).Methods("POST")
	api.HandleFunc("/unfollow", r.handler.UnfollowUserHandler).Methods("POST")

	// Health check
	api.HandleFunc("/health", r.handler.HealthCheckHandler).Methods("GET")

	// Add CORS middleware
	router.Use(corsMiddleware)

	return router
}

// corsMiddleware adds CORS headers for development
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-User-ID")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
