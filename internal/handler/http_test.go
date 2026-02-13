package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/thihxm/ebanx-home-assignment/internal/domain"
)

type MockService struct {
	BalanceFunc      func(string) (int, error)
	DepositFunc      func(string, int) (*domain.Account, error)
	WithdrawFunc     func(string, int) (*domain.Account, error)
	TransferFunc     func(string, string, int) (*domain.Account, *domain.Account, error)
	ProcessEventFunc func(domain.EventRequest) (*domain.EventResponse, error)
	ResetFunc        func() error
}

func (m *MockService) GetBalance(id string) (int, error) {
	return m.BalanceFunc(id)
}

func (m *MockService) Deposit(id string, amount int) (*domain.Account, error) {
	return m.DepositFunc(id, amount)
}

func (m *MockService) Withdraw(id string, amount int) (*domain.Account, error) {
	return m.WithdrawFunc(id, amount)
}

func (m *MockService) Transfer(originID, destinationID string, amount int) (*domain.Account, *domain.Account, error) {
	return m.TransferFunc(originID, destinationID, amount)
}

func (m *MockService) ProcessEvent(req domain.EventRequest) (*domain.EventResponse, error) {
	return m.ProcessEventFunc(req)
}

func (m *MockService) Reset() error {
	if m.ResetFunc != nil {
		return m.ResetFunc()
	}
	return nil
}

func TestGetBalance_Success(t *testing.T) {
	mockSvc := &MockService{
		BalanceFunc: func(id string) (int, error) {
			if id == "100" {
				return 20, nil
			}
			return 0, errors.New("not found")
		},
	}

	h := NewAccountHTTPHandler(mockSvc, mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/balance?account_id=100", nil)
	w := httptest.NewRecorder()

	h.handleGetBalance(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if w.Body.String() != "20" {
		t.Errorf("Expected body '20', got %q", w.Body.String())
	}
}

func TestGetBalance_NotFound(t *testing.T) {
	mockSvc := &MockService{
		BalanceFunc: func(id string) (int, error) {
			return 0, errors.New("account not found")
		},
	}

	h := NewAccountHTTPHandler(mockSvc, mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/balance?account_id=999", nil)
	w := httptest.NewRecorder()

	h.handleGetBalance(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
	if w.Body.String() != "0" {
		t.Errorf("Expected body '0', got %q", w.Body.String())
	}
}

func TestHandleEvent_Deposit(t *testing.T) {
	expectedResp := &domain.EventResponse{
		Destination: &domain.Account{ID: "100", Balance: 10},
	}

	mockSvc := &MockService{
		ProcessEventFunc: func(req domain.EventRequest) (*domain.EventResponse, error) {
			return expectedResp, nil
		},
	}

	h := NewAccountHTTPHandler(mockSvc, mockSvc)

	body := []byte(`{"type":"deposit", "destination":"100", "amount":10}`)
	req := httptest.NewRequest(http.MethodPost, "/event", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.handleEvent(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var actualResp domain.EventResponse
	json.Unmarshal(w.Body.Bytes(), &actualResp)
	if actualResp.Destination.Balance != 10 {
		t.Errorf("Expected balance 10, got %d", actualResp.Destination.Balance)
	}
}
