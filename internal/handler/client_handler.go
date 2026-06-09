package handler

import (
	"bank/internal/models"
	"bank/internal/service"
	"encoding/json"
	"net/http"
)

type ClientHandler struct {
	svc *service.ClientService
}

func NewClientHandler(s *service.ClientService) *ClientHandler {
	return &ClientHandler{svc: s}
}

func (h *ClientHandler) List(w http.ResponseWriter, r *http.Request) {
	clients, err := h.svc.FindAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(clients)
}

func (h *ClientHandler) Create(w http.ResponseWriter, r *http.Request) {
	var c models.Client
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if c.LastName == "" || c.FirstName == "" || c.Phone == "" {
		http.Error(w, "last_name, first_name and phone are required", http.StatusBadRequest)
		return
	}
	if err := h.svc.Create(&c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(c)
}

func (h *ClientHandler) Search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	var (
		clients []models.Client
		err     error
	)
	if ln := q.Get("last_name"); ln != "" {
		clients, err = h.svc.FindByLastName(ln)
	} else if ph := q.Get("phone"); ph != "" {
		clients, err = h.svc.FindByPhone(ph)
	} else {
		http.Error(w, "provide last_name or phone query param", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(clients)
}

func (h *ClientHandler) Update(w http.ResponseWriter, r *http.Request) {
	var c models.Client
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if c.ID == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	if err := h.svc.Update(&c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(c)
}

func (h *ClientHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	if err := h.svc.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ClientHandler) DeleteAllClosed(w http.ResponseWriter, r *http.Request) {
	n, err := h.svc.DeleteAllClosed()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]int64{"deleted": n})
}
