package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/quangd42/meal-planner/backend/handlers"
	"github.com/quangd42/meal-planner/backend/internal/database"

	_ "github.com/lib/pq"
)

func run() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading env file")
	}

	port := os.Getenv("PORT")
	jwtSecret := os.Getenv("JWT_SECRET")
	dbURL := os.Getenv("DB_URL")
	if port == "" || jwtSecret == "" || dbURL == "" {
		log.Fatal("missing env settings")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("error connecting to the database")
	}
	dbQueries := database.New(db)

	config := &handlers.Config{
		Port:      ":" + port,
		JWTSecret: jwtSecret,
		DB:        dbQueries,
	}

	r := chi.NewRouter()
	handlers.AddRoutes(r, config)

	fmt.Printf("listening on port %s...\n", config.Port)
	err = http.ListenAndServe(config.Port, r)
	if err != nil && err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
	}
	return nil
}
