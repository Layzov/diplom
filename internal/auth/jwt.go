package auth

import (
	"fmt"
	"time"

	"diplom/internal/apperror"
	"diplom/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTConfig holds signing parameters.
type JWTConfig struct {
	Secret     string
	AccessTTL  time.Duration
}

type Manager struct {
	secret []byte
	ttl    time.Duration
}

func NewManager(cfg JWTConfig) (*Manager, error) {
	if len(cfg.Secret) < 16 {
		return nil, fmt.Errorf("jwt secret must be at least 16 characters")
	}
	if cfg.AccessTTL <= 0 {
		cfg.AccessTTL = 24 * time.Hour
	}
	return &Manager{secret: []byte(cfg.Secret), ttl: cfg.AccessTTL}, nil
}

type accessClaims struct {
	UserID string `json:"uid"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func (m *Manager) Issue(userID uuid.UUID, email string, role models.Role) (string, time.Time, error) {
	now := time.Now().UTC()
	exp := now.Add(m.ttl)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims{
		UserID: userID.String(),
		Email:  email,
		Role:   string(role),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	})
	signed, err := token.SignedString(m.secret)
	if err != nil {
		return "", time.Time{}, err
	}
	return signed, exp, nil
}

func (m *Manager) Parse(tokenStr string) (Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &accessClaims{}, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return m.secret, nil
	})
	if err != nil {
		return Claims{}, apperror.ErrUnauthorized
	}
	claims, ok := token.Claims.(*accessClaims)
	if !ok || !token.Valid {
		return Claims{}, apperror.ErrUnauthorized
	}
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return Claims{}, apperror.ErrUnauthorized
	}
	return Claims{
		UserID: userID,
		Email:  claims.Email,
		Role:   models.Role(claims.Role),
	}, nil
}
