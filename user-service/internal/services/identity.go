package services

import (
	"context"
	"errors"
	"time"

	"github.com/HeadGardener/TaxiApp/user-service/internal/lib/auth"
	"github.com/HeadGardener/TaxiApp/user-service/internal/lib/hash"
	"github.com/HeadGardener/TaxiApp/user-service/internal/models"
	"github.com/google/uuid"
)

var (
	ErrNotActive       = errors.New("unable to get access to this account")
	ErrInvalidPassword = errors.New("invalid password")
)

type UserProcessor interface {
	Create(ctx context.Context, user *models.User) (string, error)
	GetByPhone(ctx context.Context, phone string) (models.User, error)
}

type TokenStorage interface {
	Add(ctx context.Context, userID, token string) error
	Check(ctx context.Context, userID, token string) error
	Delete(ctx context.Context, userID string) error
}

type IdentityService struct {
	userProcessor UserProcessor
	tokenStorage  TokenStorage
}

func NewIdentityService(userProcessor UserProcessor, tokenStorage TokenStorage) *IdentityService {
	return &IdentityService{
		userProcessor: userProcessor,
		tokenStorage:  tokenStorage,
	}
}

func (s *IdentityService) SignUp(ctx context.Context, user *models.User) (string, error) {
	{
		user.ID = uuid.NewString()
		user.Password = hash.GetPasswordHash(user.Password)
		user.Rating = 0.0
		user.Registration = time.Now()
		user.IsActive = true
	}

	return s.userProcessor.Create(ctx, user)
}

func (s *IdentityService) SignIn(ctx context.Context, phone, password string) (string, error) {
	user, err := s.userProcessor.GetByPhone(ctx, phone)
	if err != nil {
		return "", err
	}

	if !user.IsActive {
		return "", ErrNotActive
	}

	if !hash.CheckPassword([]byte(user.Password), password) {
		return "", ErrInvalidPassword
	}

	token, err := auth.GenerateToken(user.ID, user.Phone)
	if err != nil {
		return "", err
	}

	if err = s.tokenStorage.Add(ctx, user.ID, token); err != nil {
		return "", err
	}

	return token, nil
}

func (s *IdentityService) Check(ctx context.Context, userID, token string) error {
	return s.tokenStorage.Check(ctx, userID, token)
}

func (s *IdentityService) LogOut(ctx context.Context, userID string) error {
	return s.tokenStorage.Delete(ctx, userID)
}
