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
	account, err := s.repo.FindByID(accountID)
	if err != nil {
		return 0, err
	}
	if account == nil {
		return 0, fmt.Errorf("Account not found")
	}
	return account.Balance, nil
}

func (s *AccountService) Deposit(accountID string, amount int) (*domain.Account, error) {
	account, err := s.repo.FindByID(accountID)
	if err != nil {
		return nil, err
	}
	if account == nil {
		account = &domain.Account{
			ID:      accountID,
			Balance: 0,
		}
	}
	account.Balance += amount
	return s.repo.Upsert(account)
}

func (s *AccountService) Withdraw(accountID string, amount int) (*domain.Account, error) {
	account, err := s.repo.FindByID(accountID)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, fmt.Errorf("Account not found")
	}
	if account.Balance < amount {
		return nil, fmt.Errorf("Insufficient funds")
	}
	account.Balance -= amount
	return s.repo.Upsert(account)
}

func (s *AccountService) Transfer(originID, destinationID string, amount int) (*domain.Account, *domain.Account, error) {
	originAccount, err := s.repo.FindByID(originID)
	if err != nil {
		return nil, nil, err
	}
	if originAccount == nil {
		return nil, nil, fmt.Errorf("Origin account not found")
	}
	if originAccount.Balance < amount {
		return nil, nil, fmt.Errorf("Insufficient funds")
	}
	originAccount.Balance -= amount

	destinationAccount, err := s.repo.FindByID(destinationID)
	if err != nil {
		return nil, nil, err
	}
	if destinationAccount == nil {
		destinationAccount = &domain.Account{
			ID:      destinationID,
			Balance: 0,
		}
	}
	destinationAccount.Balance += amount

	originAccount, err = s.repo.Upsert(originAccount)
	if err != nil {
		return nil, nil, err
	}
	destinationAccount, err = s.repo.Upsert(destinationAccount)
	if err != nil {
		return nil, nil, err
	}
	return originAccount, destinationAccount, nil
}

func (s *AccountService) Reset() error {
	return s.repo.Reset()
}
