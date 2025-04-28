package server

import (
	"net/http"

	"jwt-go/internal/handler"

	"github.com/gorilla/mux"
)

func NewServer(authHandler *handler.AuthHandler) *http.Server {
	r := mux.NewRouter()

	r.HandleFunc("/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/login", authHandler.Login).Methods("POST")

	return &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
}
