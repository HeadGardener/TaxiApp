package storage

import (
	"context"
	"fmt"

	"github.com/HeadGardener/TaxiApp/user-service/internal/models"

	"github.com/jmoiron/sqlx"
)

type UserStorage struct {
	db *sqlx.DB
}

func NewUserStorage(db *sqlx.DB) *UserStorage {
	return &UserStorage{db: db}
}

func (s *UserStorage) Create(ctx context.Context, user *models.User) (string, error) {
	var createUserQuery = fmt.Sprintf(`INSERT INTO %s
    									(id, name, surname, phone, email, password_hash, rating, date, is_active)
										VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9)`, usersTable)

	if _, err := s.db.ExecContext(ctx,
		createUserQuery,
		user.ID,
		user.Name,
		user.Surname,
		user.Phone,
		user.Email,
		user.Password,
		user.Rating,
		user.Registration,
		user.IsActive); err != nil {
		return "", err
	}

	return user.ID, nil
}

func (s *UserStorage) GetByID(ctx context.Context, userID string) (models.User, error) {
	var getUserByIDQuery = fmt.Sprintf(`SELECT * FROM %s WHERE id=$1`, usersTable)

	var user models.User

	if err := s.db.GetContext(ctx, &user, getUserByIDQuery, userID); err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (s *UserStorage) GetByPhone(ctx context.Context, phone string) (models.User, error) {
	var getUserByPhoneQuery = fmt.Sprintf(`SELECT * FROM %s WHERE phone=$1`, usersTable)

	var user models.User

	if err := s.db.GetContext(ctx, &user, getUserByPhoneQuery, phone); err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (s *UserStorage) Update(ctx context.Context, userID string, userUpdate *models.User) error {
	var updateUserQuery = fmt.Sprintf(`UPDATE %s SET name=$1, surname=$2, phone=$3, email=$4 WHERE id=$5`,
		usersTable)

	if _, err := s.db.ExecContext(ctx,
		updateUserQuery,
		userUpdate.Name,
		userUpdate.Surname,
		userUpdate.Phone,
		userUpdate.Email,
		userID); err != nil {
		return err
	}

	return nil
}

func (s *UserStorage) SetInactive(ctx context.Context, userID string) error {
	var setInactiveQuery = fmt.Sprintf(`UPDATE %s SET is_active=false WHERE id=$1`, usersTable)

	if _, err := s.db.ExecContext(ctx,
		setInactiveQuery,
		userID); err != nil {
		return err
	}

	return nil
}

func (s *UserStorage) GetAll(ctx context.Context) ([]models.User, error) {
	var getAllUsersQuery = fmt.Sprintf(`SELECT * FROM %s`, usersTable)

	var users []models.User

	if err := s.db.SelectContext(ctx, &users, getAllUsersQuery); err != nil {
		return nil, err
	}

	return users, nil
}
