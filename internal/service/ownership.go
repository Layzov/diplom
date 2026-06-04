package service

import (
	"diplom/internal/apperror"

	"github.com/google/uuid"
)

func ensureOwner(ownerID, authUserID uuid.UUID) error {
	if ownerID != authUserID {
		return apperror.ErrForbidden
	}
	return nil
}
