package services

import (
	"errors"

	"github.com/jackc/pgx/v5"
)

var ErrResourceNotFound = errors.New("resource not found")

func checkErrNoRows(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrResourceNotFound
	}
	return err
}
