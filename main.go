package main

import (
	"fmt"
	"log"
	"net/http"

	"uala-challenge/internal/application/usecases"
	httpInterface "uala-challenge/internal/interfaces/http"
	"uala-challenge/internal/infrastructure/storage"
)

func main() {
	// Initialize infrastructure layer
	inMemoryStorage := storage.NewInMemoryRepository()
	userRepo := storage.NewUserRepository(inMemoryStorage)
	tweetRepo := storage.NewTweetRepository(inMemoryStorage)
	followRepo := storage.NewFollowRepository(inMemoryStorage)

	// Initialize application layer (use cases)
	tweetUseCase := usecases.NewTweetUseCase(tweetRepo, userRepo)
	followUseCase := usecases.NewFollowUseCase(followRepo, tweetRepo)

	// Initialize interface layer (HTTP handlers)
	handler := httpInterface.NewHandler(tweetUseCase, followUseCase)
	router := httpInterface.NewRouter(handler)

	// Setup routes
	httpRouter := router.SetupRoutes()

	// Start server
	port := ":8080"
	fmt.Printf("Server starting on port %s\n", port)
	fmt.Println("Available endpoints:")
	fmt.Println("  POST   /api/v1/tweets         - Create a tweet")
	fmt.Println("  GET    /api/v1/timeline       - Get user timeline")
	fmt.Println("  GET    /api/v1/users/tweets   - Get user tweets")
	fmt.Println("  POST   /api/v1/follow         - Follow a user")
	fmt.Println("  POST   /api/v1/unfollow       - Unfollow a user")
	fmt.Println("  GET    /api/v1/health         - Health check")
	fmt.Println("\nNote: Include X-User-ID header for user identification")

	log.Fatal(http.ListenAndServe(port, httpRouter))
}