package service

import (
	"testing"

	"github.com/thihxm/ebanx-home-assignment/internal/domain"
	"github.com/thihxm/ebanx-home-assignment/internal/repository"
)

func TestHandleDepositEvent(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	accountService := NewAccountService(repo)
	eventService := NewEventService(accountService)

	_, err := eventService.HandleEvent(domain.EventRequest{
		Type:   "deposit",
		Origin: "123",
		Amount: 100,
	})
	if err != nil {
		t.Errorf("Expected no error depositing: %v", err)
	}
}

func TestHandleWithdrawEvent(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	accountService := NewAccountService(repo)
	eventService := NewEventService(accountService)

	accountService.Deposit("123", 100)

	_, err := eventService.HandleEvent(domain.EventRequest{
		Type:   "withdraw",
		Origin: "123",
		Amount: 100,
	})
	if err != nil {
		t.Errorf("Expected no error withdrawing: %v", err)
	}
}

func TestHandleTransferEvent(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	accountService := NewAccountService(repo)
	eventService := NewEventService(accountService)

	accountService.Deposit("123", 100)
	accountService.Deposit("456", 0)

	_, err := eventService.HandleEvent(domain.EventRequest{
		Type:        "transfer",
		Origin:      "123",
		Destination: "456",
		Amount:      100,
	})
	if err != nil {
		t.Errorf("Expected no error transferring: %v", err)
	}
}
