package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/quangd42/meal-planner/internal/database"
	"github.com/quangd42/meal-planner/internal/handlers"
	"github.com/quangd42/meal-planner/internal/services"

	_ "github.com/lib/pq"
)

func run() error {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Fatal("error loading env")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}

	db, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()
	store := database.NewStore(db)

	us := services.NewUserService(store)
	ts := services.NewTokenService(store)

	r := chi.NewRouter()
	handlers.AddRoutes(r, us, ts)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Printf("listening on port %s...\n", port)
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
	}
	return nil
}
