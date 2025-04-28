package server

import (
	"jwt-go/internal/handler"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(authHandler *handler.AuthHandler, productHandler *handler.ProductHandler) *Server {
	router := mux.NewRouter()

	// Public
	router.HandleFunc("/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/login", authHandler.Login).Methods("POST")

	// Protected
	productsRouter := router.PathPrefix("/products").Subrouter()
	productsRouter.Use(JWTMiddleware("your-super-secret"))
	productsRouter.HandleFunc("", productHandler.ListProducts).Methods("GET")

	return &Server{
		httpServer: &http.Server{
			Addr:    ":8080",
			Handler: router,
		},
	}
}

func (s *Server) ListenAndServe() error {
	return s.httpServer.ListenAndServe()
}
