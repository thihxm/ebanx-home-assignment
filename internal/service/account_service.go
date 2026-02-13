package service

import (
	"fmt"

	"github.com/thihxm/ebanx-home-assignment/internal/domain"
)

type AccountService struct {
	repo domain.AccountRepository
}

func NewAccountService(repo domain.AccountRepository) *AccountService {
	return &AccountService{
		repo: repo,
	}
}

func (s *AccountService) GetBalance(accountID string) (int, error) {
	return 0, fmt.Errorf("Not implemented")
}

func (s *AccountService) Deposit(accountID string, amount int) (*domain.Account, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s *AccountService) Withdraw(accountID string, amount int) (*domain.Account, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s *AccountService) Transfer(originID, destinationID string, amount int) (*domain.Account, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s *AccountService) Reset() error {
	return fmt.Errorf("Not implemented")
}
