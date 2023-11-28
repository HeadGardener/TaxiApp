package models

import "time"

type User struct {
	ID           string    `db:"id"`
	Name         string    `db:"name"`
	Surname      string    `db:"surname"`
	Phone        string    `db:"phone"`
	Email        string    `db:"email"`
	Password     string    `db:"password_hash"`
	Rating       float32   `db:"rating"`
	Registration time.Time `db:"date"`
	IsActive     bool      `db:"is_active"`
}

type UserProfile struct {
	Name    string
	Surname string
	Phone   string
	Email   string
	Rating  float32
}
