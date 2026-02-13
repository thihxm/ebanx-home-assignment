package handler

import (
	"github.com/thihxm/ebanx-home-assignment/internal/domain"
)

type HTTPHandler struct {
	accountUseCase domain.AccountUseCase
}

func NewHTTPHandler(accountUseCase domain.AccountUseCase) *HTTPHandler {
	return &HTTPHandler{
		accountUseCase: accountUseCase,
	}
}
