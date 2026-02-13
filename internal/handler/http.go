package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/thihxm/ebanx-home-assignment/internal/domain"
)

var uni *ut.UniversalTranslator

type HTTPHandler struct {
	accountService domain.AccountService
	eventService   domain.EventService
	validate       *validator.Validate
}

func NewAccountHTTPHandler(accountService domain.AccountService, eventService domain.EventService) *HTTPHandler {
	en := en.New()
	uni = ut.New(en, en)
	validate := validator.New()
	enTrans, _ := uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(validate, enTrans)

	return &HTTPHandler{
		accountService: accountService,
		eventService:   eventService,
		validate:       validate,
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
	err := h.validate.Struct(req)
	if err != nil {
		var invalidValidationError *validator.InvalidValidationError
		if errors.As(err, &invalidValidationError) {
			fmt.Println(err)
			return
		}

		var errs validator.ValidationErrors
		var httpErrors []validator.ValidationErrorsTranslations
		trans, _ := uni.GetTranslator(strings.Replace(r.Header.Get("Accept-Language"), "-", "_", -1))

		if errors.As(err, &errs) {
			httpErrors = append(httpErrors, errs.Translate(trans))
		}
		r, _ := json.Marshal(httpErrors)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(r)
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
		fmt.Fprintf(w, "missing account_id")
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
