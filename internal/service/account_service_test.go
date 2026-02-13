package service

import (
	"testing"

	"github.com/thihxm/ebanx-home-assignment/internal/repository"
)

func TestGetNonExistentAccount(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	service := NewAccountService(repo)

	_, err := service.GetBalance("123")
	if err == nil {
		t.Errorf("Expected error getting balance")
	}
}

func TestGetBalance(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	service := NewAccountService(repo)

	_, err := service.Deposit("123", 100)
	if err != nil {
		t.Errorf("Expected no error depositing: %v", err)
	}

	balance, err := service.GetBalance("123")
	if err != nil {
		t.Errorf("Expected no error getting balance: %v", err)
	}
	if balance != 100 {
		t.Errorf("Expected balance to be 100, got %d", balance)
	}
}

func TestDeposit(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	service := NewAccountService(repo)

	account, err := service.Deposit("123", 100)
	if err != nil {
		t.Errorf("Expected no error depositing: %v", err)
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

func TestWithdrawNonExistentAccount(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	service := NewAccountService(repo)

	_, err := service.Withdraw("123", 100)
	if err == nil {
		t.Errorf("Expected error withdrawing")
	}
}

func TestWithdraw(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	service := NewAccountService(repo)

	_, err := service.Deposit("123", 100)
	if err != nil {
		t.Errorf("Expected no error depositing: %v", err)
	}

	account, err := service.Withdraw("123", 100)
	if err != nil {
		t.Errorf("Expected no error withdrawing: %v", err)
	}
	if account == nil {
		t.Errorf("Expected non-nil account")
	}
	if account.ID != "123" {
		t.Errorf("Expected account ID to be 123")
	}
	if account.Balance != 0 {
		t.Errorf("Expected account balance to be 0")
	}
}

func TestWithdrawInsufficientFunds(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	service := NewAccountService(repo)

	_, err := service.Deposit("123", 100)
	if err != nil {
		t.Errorf("Expected no error depositing: %v", err)
	}

	_, err = service.Withdraw("123", 200)
	if err == nil {
		t.Errorf("Expected error withdrawing")
	}
}

func TestTransferNonExistentAccounts(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	service := NewAccountService(repo)

	_, err := service.Transfer("123", "456", 100)
	if err == nil {
		t.Errorf("Expected error transferring")
	}
}

func TestTransferNonExistentOriginAccount(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	service := NewAccountService(repo)

	_, err := service.Deposit("456", 100)
	if err != nil {
		t.Errorf("Expected no error depositing: %v", err)
	}

	_, err = service.Transfer("123", "456", 100)
	if err == nil {
		t.Errorf("Expected error transferring")
	}
}

func TestTransferNonExistentDestinationAccount(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	service := NewAccountService(repo)

	_, err := service.Deposit("123", 100)
	if err != nil {
		t.Errorf("Expected no error depositing: %v", err)
	}

	_, err = service.Transfer("123", "456", 100)
	if err == nil {
		t.Errorf("Expected error transferring")
	}
}

func TestTransfer(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	service := NewAccountService(repo)

	_, err := service.Deposit("123", 100)
	if err != nil {
		t.Errorf("Expected no error depositing: %v", err)
	}

	_, err = service.Deposit("456", 100)
	if err != nil {
		t.Errorf("Expected no error depositing: %v", err)
	}

	account, err := service.Transfer("123", "456", 100)
	if err != nil {
		t.Errorf("Expected no error transferring: %v", err)
	}
	if account == nil {
		t.Errorf("Expected non-nil account")
	}
	if account.ID != "123" {
		t.Errorf("Expected account ID to be 123")
	}
	if account.Balance != 0 {
		t.Errorf("Expected account balance to be 0")
	}

	destinationAccount, err := service.GetBalance("456")
	if err != nil {
		t.Errorf("Expected no error getting balance: %v", err)
	}
	if destinationAccount != 200 {
		t.Errorf("Expected account balance to be 200")
	}
}

func TestTransferInsufficientFunds(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	service := NewAccountService(repo)

	_, err := service.Deposit("123", 100)
	if err != nil {
		t.Errorf("Expected no error depositing: %v", err)
	}

	_, err = service.Deposit("456", 0)
	if err != nil {
		t.Errorf("Expected no error depositing: %v", err)
	}

	_, err = service.Transfer("123", "456", 200)
	if err == nil {
		t.Errorf("Expected error transferring")
	}
}

func TestReset(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	service := NewAccountService(repo)

	_, err := service.Deposit("123", 100)
	if err != nil {
		t.Errorf("Expected no error depositing: %v", err)
	}

	err = service.Reset()
	if err != nil {
		t.Errorf("Expected no error resetting: %v", err)
	}

	_, err = service.GetBalance("123")
	if err == nil {
		t.Errorf("Expected error getting balance")
	}
}
