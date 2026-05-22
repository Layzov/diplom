package service

import (
	"strings"

	"diplom/internal/apperror"
	"diplom/internal/dto"
	"diplom/internal/models"
	"diplom/internal/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(req dto.CreateUserRequest) (*dto.UserResponse, error) {
	if err := validateEmail(req.Email); err != nil {
		return nil, err
	}
	if err := validateNonEmpty(req.Name, "name"); err != nil {
		return nil, err
	}
	if len(req.Password) < 8 {
		return nil, apperror.Validation("password must be at least 8 characters")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	role := models.UserRoleUser
	if req.Role != nil {
		if *req.Role != models.UserRoleAdmin && *req.Role != models.UserRoleUser {
			return nil, apperror.Validation("invalid role")
		}
		role = *req.Role
	}

	user := &models.User{
		Email:        strings.TrimSpace(req.Email),
		Name:         strings.TrimSpace(req.Name),
		Role:         role,
		PasswordHash: string(hash),
	}
	if req.Settings != nil {
		user.Settings = models.UserSettings{
			Timezone:             req.Settings.Timezone,
			Language:             req.Settings.Language,
			Theme:                req.Settings.Theme,
			NotificationsEnabled: req.Settings.NotificationsEnabled,
		}
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	resp := dto.UserFromModel(user)
	return &resp, nil
}

func (s *UserService) List() (dto.UserListResponse, error) {
	users, err := s.repo.List()
	if err != nil {
		return dto.UserListResponse{}, err
	}
	items := make([]dto.UserResponse, 0, len(users))
	for i := range users {
		items = append(items, dto.UserFromModel(&users[i]))
	}
	return dto.UserListResponse{Items: items}, nil
}

func (s *UserService) Get(id uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	resp := dto.UserFromModel(user)
	return &resp, nil
}

func (s *UserService) Update(id uuid.UUID, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Email != nil {
		if err := validateEmail(*req.Email); err != nil {
			return nil, err
		}
		user.Email = strings.TrimSpace(*req.Email)
	}
	if req.Name != nil {
		if err := validateNonEmpty(*req.Name, "name"); err != nil {
			return nil, err
		}
		user.Name = strings.TrimSpace(*req.Name)
	}
	if req.Password != nil {
		if len(*req.Password) < 8 {
			return nil, apperror.Validation("password must be at least 8 characters")
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.PasswordHash = string(hash)
	}
	if req.Role != nil {
		if *req.Role != models.UserRoleAdmin && *req.Role != models.UserRoleUser {
			return nil, apperror.Validation("invalid role")
		}
		user.Role = *req.Role
	}
	if req.Settings != nil {
		if req.Settings.Timezone != "" {
			user.Settings.Timezone = req.Settings.Timezone
		}
		if req.Settings.Language != "" {
			user.Settings.Language = req.Settings.Language
		}
		if req.Settings.Theme != "" {
			user.Settings.Theme = req.Settings.Theme
		}
		user.Settings.NotificationsEnabled = req.Settings.NotificationsEnabled
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}
	resp := dto.UserFromModel(user)
	return &resp, nil
}

func (s *UserService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}
