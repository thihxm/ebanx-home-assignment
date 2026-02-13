package main

import (
	"log"

	"github.com/thihxm/ebanx-home-assignment/internal/handler"
	"github.com/thihxm/ebanx-home-assignment/internal/repository"
	"github.com/thihxm/ebanx-home-assignment/internal/service"
)

func main() {
	repo := repository.NewInMemoryRepository()
	accountService := service.NewAccountService(repo)
	eventService := service.NewEventService(accountService)

	httpHandler := handler.NewAccountHTTPHandler(accountService, eventService)

	if err := httpHandler.Serve(":8080"); err != nil {
		log.Fatalf("Error serving HTTP server: %v", err)
	}
}
