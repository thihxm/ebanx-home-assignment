package repository

import (
	"fmt"

	"github.com/thihxm/ebanx-home-assignment/internal/domain"
)

type InMemoryRepository struct {
	accounts map[string]*domain.Account
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		accounts: make(map[string]*domain.Account),
	}
}

func (r *InMemoryRepository) FindByID(id string) (*domain.Account, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (r *InMemoryRepository) Upsert(account *domain.Account) (*domain.Account, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (r *InMemoryRepository) Reset() error {
	return fmt.Errorf("Not implemented")
}
