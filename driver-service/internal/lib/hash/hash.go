package hash

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	cost = 8
)

func GetPasswordHash(password string) string {
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), cost)

	return string(passwordHash)
}

func CheckPassword(passHash []byte, password string) bool {
	err := bcrypt.CompareHashAndPassword(passHash, []byte(password))

	return err == nil
}
