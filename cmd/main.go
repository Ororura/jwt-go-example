package main

import (
	"log"

	"jwt-go/internal/handler"
	"jwt-go/internal/repository"
	"jwt-go/internal/server"
	"jwt-go/internal/usecase"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sqlx.Open("sqlite3", "./auth.db")
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer db.Close()

	if err := server.RunMigrations(db); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	userRepo := repository.NewUserSQLiteRepository(db)
	authUC := usecase.NewAuthUsecase(userRepo, "your-super-secret")
	authHandler := handler.NewAuthHandler(authUC)

	productRepo := repository.NewProductSQLiteRepository(db)
	productUC := usecase.NewProductUsecase(productRepo)
	productHandler := handler.NewProductHandler(productUC)

	srv := server.NewServer(authHandler, productHandler)

	log.Println("Starting server on :8080")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
