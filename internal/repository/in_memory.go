package repository

import (
	"fmt"
	"sync"

	"github.com/thihxm/ebanx-home-assignment/internal/domain"
)

type InMemoryRepository struct {
	accounts map[string]*domain.Account
	mu       sync.RWMutex
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		accounts: make(map[string]*domain.Account),
	}
}

func (r *InMemoryRepository) FindByID(id string) (*domain.Account, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	account, ok := r.accounts[id]
	if !ok {
		return nil, nil
	}
	return account, nil
}

func (r *InMemoryRepository) Upsert(account *domain.Account) (*domain.Account, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.accounts[account.ID] = account

	return account, nil
}

func (r *InMemoryRepository) Reset() error {
	return fmt.Errorf("Not implemented")
}
