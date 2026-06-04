package service

import (
	"diplom/internal/auth"
	"diplom/internal/dto"
)

type AuthService struct {
	users *UserService
	jwt   *auth.Manager
}

func NewAuthService(users *UserService, jwt *auth.Manager) *AuthService {
	return &AuthService{users: users, jwt: jwt}
}

func (s *AuthService) Register(req dto.RegisterRequest) (*dto.AuthResponse, error) {
	create := dto.CreateUserRequest{
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password,
		Settings: req.Settings,
	}
	user, err := s.users.Create(create)
	if err != nil {
		return nil, err
	}
	return s.tokenResponse(user)
}

func (s *AuthService) Login(req dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.users.Authenticate(req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	return s.tokenResponse(user)
}

func (s *AuthService) tokenResponse(user *dto.UserResponse) (*dto.AuthResponse, error) {
	token, exp, err := s.jwt.Issue(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}
	return &dto.AuthResponse{
		AccessToken: token,
		ExpiresAt:   exp,
		User:        *user,
	}, nil
}
