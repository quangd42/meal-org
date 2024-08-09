package services

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrResourceNotFound = errors.New("resource not found")
	ErrDBConstraint     = errors.New("data constraint error")
	ErrUniqueValue      = errors.New("unique value error")
)

func checkErrNoRows(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrResourceNotFound
	}
	return err
}

func checkErrDBConstraint(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code[0:2] == "23" {
		return ErrDBConstraint
	}
	return err
}

func customDBErr(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrResourceNotFound
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code[0:2] == "23" {
		return ErrDBConstraint
	}
	return err
}
