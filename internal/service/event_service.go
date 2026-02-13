package service

import (
	"fmt"

	"github.com/thihxm/ebanx-home-assignment/internal/domain"
)

type EventService struct {
	accountService domain.AccountService
}

func NewEventService(accountService domain.AccountService) *EventService {
	return &EventService{
		accountService: accountService,
	}
}

func (s *EventService) ProcessEvent(event domain.EventRequest) (*domain.EventResponse, error) {
	switch event.Type {
	case "deposit":
		account, err := s.accountService.Deposit(event.Destination, event.Amount)
		if err != nil {
			return nil, err
		}
		return &domain.EventResponse{
			Destination: account,
		}, nil
	case "withdraw":
		account, err := s.accountService.Withdraw(event.Origin, event.Amount)
		if err != nil {
			return nil, err
		}
		return &domain.EventResponse{
			Origin: account,
		}, nil
	case "transfer":
		originAccount, destinationAccount, err := s.accountService.Transfer(event.Origin, event.Destination, event.Amount)
		if err != nil {
			return nil, err
		}
		return &domain.EventResponse{
			Origin:      originAccount,
			Destination: destinationAccount,
		}, nil
	default:
		return nil, fmt.Errorf("invalid event type: %s", event.Type)
	}
}
