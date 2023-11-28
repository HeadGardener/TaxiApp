package services

import (
	"context"

	"github.com/HeadGardener/TaxiApp/user-service/internal/models"
)

type UserStorage interface {
	GetByID(ctx context.Context, userID string) (models.User, error)
	Update(ctx context.Context, userID string, userUpdate *models.User) error
	SetInactive(ctx context.Context, userID string) error
}

type UserService struct {
	userStorage UserStorage
}

func NewUserService(userStorage UserStorage) *UserService {
	return &UserService{
		userStorage: userStorage,
	}
}

func (s *UserService) GetProfile(ctx context.Context, userID string) (models.UserProfile, error) {
	user, err := s.userStorage.GetByID(ctx, userID)
	if err != nil {
		return models.UserProfile{}, err
	}

	profile := models.UserProfile{
		Name:    user.Name,
		Surname: user.Surname,
		Phone:   user.Phone,
		Email:   user.Email,
		Rating:  user.Rating,
	}

	return profile, nil
}

func (s *UserService) Update(ctx context.Context, userID string, userUpdate *models.User) error {
	user, err := s.userStorage.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if userUpdate.Name == "" {
		userUpdate.Name = user.Name
	}

	if userUpdate.Surname == "" {
		userUpdate.Surname = user.Surname
	}

	if userUpdate.Phone == "" {
		userUpdate.Phone = user.Phone
	}

	if userUpdate.Email == "" {
		userUpdate.Email = user.Email
	}

	return s.userStorage.Update(ctx, userID, userUpdate)
}

func (s *UserService) SetInactive(ctx context.Context, userID string) error {
	return s.userStorage.SetInactive(ctx, userID)
}
