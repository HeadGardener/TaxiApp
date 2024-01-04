package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	TokenTTL  = 24 * time.Hour          // placeholder
	secretKey = "qazwsxedcrfvtgbyhnujm" // placeholder
)

type UserAttributes struct {
	ID    string `json:"id"`
	Phone string `json:"phone"`
}

type tokenClaims struct {
	jwt.RegisteredClaims
	UserID string `json:"id"`
	Phone  string `json:"phone"`
}

func GenerateToken(userID, phone string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		userID,
		phone,
	})

	return token.SignedString([]byte(secretKey))
}

func ParseToken(accessToken string) (UserAttributes, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(secretKey), nil
	})

	if err != nil {
		return UserAttributes{}, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return UserAttributes{}, errors.New("token claims are not of type *tokenClaims")
	}

	return UserAttributes{
		ID:    claims.UserID,
		Phone: claims.Phone,
	}, nil
}
