package handler

import (
	"testing"

	"github.com/thihxm/ebanx-home-assignment/internal/repository"
	"github.com/thihxm/ebanx-home-assignment/internal/service"
)

func TestCreateHTTPHandler(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	service := service.NewAccountService(repo)
	httpHandler := NewHTTPHandler(service)
	if httpHandler == nil {
		t.Errorf("Expected HTTPHandler, got nil")
	}
}
