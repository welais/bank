package main

import (
	"bank/internal/db"
	"bank/internal/handler"
	"bank/internal/repository"
	"bank/internal/service"
	"log"
	"net/http"
)

func main() {
	database := db.Connect()

	clientRepo := repository.NewClientRepo(database)

	clientSvc := service.NewClientService(clientRepo)

	clientH := handler.NewClientHandler(clientSvc)

	// Static files
	http.Handle("/", http.FileServer(http.Dir("static")))

	http.HandleFunc("/clients", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			clientH.List(w, r)
		case http.MethodPost:
			clientH.Create(w, r)
		case http.MethodPut:
			clientH.Update(w, r)
		case http.MethodDelete:
			clientH.Delete(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/clients/search", clientH.Search)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", corsMiddleware(http.DefaultServeMux)))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if len(r.URL.Path) > 1 && r.URL.Path != "/" {
			w.Header().Set("Content-Type", "application/json")
		}
		next.ServeHTTP(w, r)
	})
}
