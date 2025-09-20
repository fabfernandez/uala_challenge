package storage

import (
	"context"

	"uala-challenge/internal/domain"
)

// UserRepository implements domain.UserRepository
type UserRepository struct {
	storage *InMemoryRepository
}

// NewUserRepository creates a new user repository
func NewUserRepository(storage *InMemoryRepository) *UserRepository {
	return &UserRepository{
		storage: storage,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	return r.storage.CreateUser(ctx, user)
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return r.storage.GetUser(ctx, id)
}
