package repository

import (
	"time"

	"diplom/internal/dto"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CalendarRepository struct {
	db *gorm.DB
}

func NewCalendarRepository(db *gorm.DB) *CalendarRepository {
	return &CalendarRepository{db: db}
}

// List reads planned repetitions from the calendar view.
func (r *CalendarRepository) List(userID uuid.UUID, from, to *time.Time) ([]dto.CalendarEntry, error) {
	q := r.db.Table("calendar").Where("user_id = ?", userID)
	if from != nil {
		q = q.Where("due_at >= ?", *from)
	}
	if to != nil {
		q = q.Where("due_at <= ?", *to)
	}
	var items []dto.CalendarEntry
	err := q.Order("due_at ASC").Find(&items).Error
	return items, mapError(err)
}
