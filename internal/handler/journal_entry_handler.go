package handler

import (
	"bank/internal/models"
	"bank/internal/service"
	"encoding/json"
	"net/http"
	"time"
)

type JournalEntryHandler struct {
	svc *service.JournalEntryService
}

func NewJournalEntryHandler(s *service.JournalEntryService) *JournalEntryHandler {
	return &JournalEntryHandler{svc: s}
}

func (h *JournalEntryHandler) List(w http.ResponseWriter, r *http.Request) {
	list, err := h.svc.FindAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(list)
}

func (h *JournalEntryHandler) Search(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "invalid date (YYYY-MM-DD)", http.StatusBadRequest)
		return
	}
	list, err := h.svc.FindByDate(date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(list)
}

func (h *JournalEntryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var e models.JournalEntry
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if e.EntryNumber == "" || e.DebitAccountID == "" || e.CreditAccountID == "" || e.Amount <= 0 {
		http.Error(w, "entry_number, debit_account_id, credit_account_id and amount>0 are required", http.StatusBadRequest)
		return
	}
	if err := h.svc.Create(&e); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(e)
}

func (h *JournalEntryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	if err := h.svc.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
