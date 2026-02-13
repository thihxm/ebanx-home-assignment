package repository

import (
	"testing"

	"github.com/thihxm/ebanx-home-assignment/internal/domain"
)

func TestNonExistentAccount(t *testing.T) {
	repo := NewInMemoryRepository()

	account, err := repo.FindByID("non-existent")

	if err != nil {
		t.Errorf("Expected error: %v", err)
	}

	if account != nil {
		t.Errorf("Expected nil account")
	}
}

func TestUpsertAccount(t *testing.T) {
	repo := NewInMemoryRepository()

	account, err := repo.Upsert(&domain.Account{
		ID:      "123",
		Balance: 100,
	})

	if err != nil {
		t.Errorf("Expected no error creating account: %v", err)
	}

	if account == nil {
		t.Errorf("Expected non-nil account")
	}

	if account.ID != "123" {
		t.Errorf("Expected account ID to be 123")
	}

	if account.Balance != 100 {
		t.Errorf("Expected account balance to be 100")
	}
}

func TestExistentAccount(t *testing.T) {
	repo := NewInMemoryRepository()

	_, err := repo.Upsert(&domain.Account{
		ID:      "123",
		Balance: 100,
	})

	if err != nil {
		t.Errorf("Expected no error creating account: %v", err)
	}

	account, err := repo.FindByID("123")

	if err != nil {
		t.Errorf("Expected no error finding account: %v", err)
	}

	if account == nil {
		t.Errorf("Expected account to exist")
	}

	if account.ID != "123" {
		t.Errorf("Expected account ID to be 123")
	}

	if account.Balance != 100 {
		t.Errorf("Expected account balance to be 100")
	}
}

func TestUpdateAccount(t *testing.T) {
	repo := NewInMemoryRepository()

	account, err := repo.Upsert(&domain.Account{
		ID:      "123",
		Balance: 100,
	})

	if err != nil {
		t.Errorf("Expected no error creating account: %v", err)
	}
	account.Balance = 200

	account, err = repo.Upsert(account)

	if err != nil {
		t.Errorf("Expected no error updating account: %v", err)
	}

	if account == nil {
		t.Errorf("Expected non-nil account")
	}

	if account.ID != "123" {
		t.Errorf("Expected account ID to be 123")
	}

	if account.Balance != 200 {
		t.Errorf("Expected account balance to be 200")
	}
}

func TestReset(t *testing.T) {
	repo := NewInMemoryRepository()

	_, err := repo.Upsert(&domain.Account{
		ID:      "123",
		Balance: 100,
	})

	if err != nil {
		t.Errorf("Expected no error creating account: %v", err)
	}

	err = repo.Reset()

	if err != nil {
		t.Errorf("Expected no error resetting accounts: %v", err)
	}

	account, err := repo.FindByID("123")

	if err != nil {
		t.Errorf("Expected no error finding account: %v", err)
	}

	if account != nil {
		t.Errorf("Expected nil account")
	}
}
