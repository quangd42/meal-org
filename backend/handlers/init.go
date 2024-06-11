package handlers

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/quangd42/meal-planner/backend/internal/database"
)

var DB *database.Queries

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading env file")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("missing env settings")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("error connecting to the database")
	}

	DB = database.New(db)
}
