package repository

import (
	"diplom/internal/apperror"
	"diplom/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	return mapError(r.db.Create(user).Error)
}

func (r *UserRepository) List() ([]models.User, error) {
	var users []models.User
	err := r.db.Order("created_at DESC").Find(&users).Error
	return users, mapError(err)
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "email = ?", email).Error
	if err != nil {
		return nil, mapError(err)
	}
	return &user, nil
}

func (r *UserRepository) GetByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, mapError(err)
	}
	return &user, nil
}

func (r *UserRepository) Update(user *models.User) error {
	return mapError(r.db.Save(user).Error)
}

func (r *UserRepository) Delete(id uuid.UUID) error {
	res := r.db.Delete(&models.User{}, "id = ?", id)
	if res.Error != nil {
		return mapError(res.Error)
	}
	if res.RowsAffected == 0 {
		return apperror.ErrNotFound
	}
	return nil
}

func (r *UserRepository) Exists(id uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("id = ?", id).Count(&count).Error
	return count > 0, mapError(err)
}
