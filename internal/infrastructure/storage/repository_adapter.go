package storage

import ()

// RepositoryAdapter adapts InMemoryRepository to implement all domain repository interfaces
type RepositoryAdapter struct {
	*InMemoryRepository
}

// NewRepositoryAdapter creates a new repository adapter
func NewRepositoryAdapter() *RepositoryAdapter {
	return &RepositoryAdapter{
		InMemoryRepository: NewInMemoryRepository(),
	}
}
