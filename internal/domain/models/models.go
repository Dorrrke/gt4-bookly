package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UID          uuid.UUID `json:"uuid"`
	Name         string    `json:"name,omitempty" validate:"required"`
	Email        string    `json:"email" validate:"required,email"`
	Passoword    string    `json:"pass" validate:"required,min=8"`
	Age          int       `json:"age,omitempty" validate:"gte=14"`
	RegisteredAt time.Time `json:"registered_at,omitempty"`
}

type UserLogin struct {
	Email     string `json:"email" validate:"required,email"`
	Passoword string `json:"pass" validate:"required,min=8"`
}

type Book struct {
	BID         uuid.UUID `json:"bid"`
	Lable       string    `json:"lable" validate:"required"`
	Author      string    `json:"author" validate:"required"`
	Description string    `json:"desc" validate:"required"`
	WritedAt    time.Time `json:"writed_at" validate:"required"`
}
