package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/thihxm/ebanx-home-assignment/internal/domain"
)

type HTTPHandler struct {
	accountService domain.AccountService
	eventService   domain.EventService
}

func NewAccountHTTPHandler(accountService domain.AccountService, eventService domain.EventService) *HTTPHandler {
	return &HTTPHandler{
		accountService: accountService,
		eventService:   eventService,
	}
}

func (h *HTTPHandler) registerRoutes(mux *http.ServeMux) error {
	mux.HandleFunc("/reset", h.handleReset)
	mux.HandleFunc("/event", h.handleEvent)
	mux.HandleFunc("/balance", h.handleGetBalance)
	return nil
}

func (h *HTTPHandler) handleReset(w http.ResponseWriter, r *http.Request) {
	if err := h.accountService.Reset(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func (h *HTTPHandler) handleEvent(w http.ResponseWriter, r *http.Request) {
	var req domain.EventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	resp, err := h.eventService.ProcessEvent(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "0")
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *HTTPHandler) handleGetBalance(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("account_id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	balance, err := h.accountService.GetBalance(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "0")
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%d", balance)
}

func (h *HTTPHandler) Serve(addr string) error {
	mux := http.NewServeMux()
	h.registerRoutes(mux)
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	log.Printf("Server started on %s", addr)

	return server.ListenAndServe()
}
