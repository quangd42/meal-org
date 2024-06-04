package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/quangd42/meal-planner/backend/handlers"
)

type Config struct {
	Port string
}

func run() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loanding env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("port is missing in env")
	}

	config := &Config{
		Port: ":" + port,
	}

	r := chi.NewRouter()
	handlers.AddRoutes(r)

	fmt.Printf("listening on port %s...\n", config.Port)
	err = http.ListenAndServe(config.Port, r)
	if err != nil && err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
	}
	return nil
}
