package database

import "github.com/jackc/pgx/v5/pgxpool"

type Store struct {
	Q  *Queries
	DB *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{
		DB: db,
		Q:  New(db),
	}
}
