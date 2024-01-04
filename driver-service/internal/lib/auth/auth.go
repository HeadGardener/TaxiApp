package auth

import (
	"errors"
	"time"

	"github.com/HeadGardener/TaxiApp/driver-service/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

const (
	TokenTTL  = 24 * time.Hour
	secretKey = "qazwsxedcrfvtgbyhnujm"
)

type tokenClaims struct {
	jwt.RegisteredClaims
	DriverID string          `json:"id"`
	Phone    string          `json:"phone"`
	TaxiType models.TaxiType `json:"taxi_type"`
}

type DriverAttributes struct {
	ID       string          `json:"id"`
	Phone    string          `json:"phone"`
	TaxiType models.TaxiType `json:"taxi_type"`
}

func GenerateToken(driverID, phone string, taxiType models.TaxiType) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		driverID,
		phone,
		taxiType,
	})

	return token.SignedString([]byte(secretKey))
}

// ParseToken will be used in middleware
func ParseToken(accessToken string) (DriverAttributes, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(secretKey), nil
	})

	if err != nil {
		return DriverAttributes{}, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return DriverAttributes{}, errors.New("token claims are not of type *tokenClaims")
	}

	return DriverAttributes{
		ID:       claims.DriverID,
		Phone:    claims.Phone,
		TaxiType: claims.TaxiType,
	}, nil
}
