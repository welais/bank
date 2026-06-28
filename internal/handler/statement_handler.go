package handler

import (
	"bank/internal/service"
	"encoding/json"
	"net/http"
	"time"
)

type StatementHandler struct {
	svc *service.StatementService
}

func NewStatementHandler(s *service.StatementService) *StatementHandler {
	return &StatementHandler{svc: s}
}

func (h *StatementHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	from, to, err := parseDateRange(q.Get("from"), q.Get("to"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if accountID := q.Get("account_id"); accountID != "" {
		list, err := h.svc.GetByAccount(accountID, from, to)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(list)
		return
	}

	list, err := h.svc.GetAll(from, to)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(list)
}

func parseDateRange(fromStr, toStr string) (time.Time, time.Time, error) {
	from := time.Time{}
	to := time.Now()

	if fromStr != "" {
		var err error
		from, err = time.Parse("2006-01-02", fromStr)
		if err != nil {
			return from, to, err
		}
	}
	if toStr != "" {
		var err error
		to, err = time.Parse("2006-01-02", toStr)
		if err != nil {
			return from, to, err
		}
	}
	if from.IsZero() {
		from = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	}
	return from, to, nil
}
