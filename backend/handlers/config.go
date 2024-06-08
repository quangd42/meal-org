package handlers

import "github.com/quangd42/meal-planner/backend/internal/database"

type Config struct {
	Port      string
	JWTSecret string
	DB        *database.Queries
}
