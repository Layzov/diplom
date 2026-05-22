package repository

import (
	"errors"

	"diplom/internal/apperror"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

func mapError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperror.ErrNotFound
	}
	var pg *pgconn.PgError
	if errors.As(err, &pg) && pg.Code == "23505" {
		return apperror.ErrConflict
	}
	return err
}
