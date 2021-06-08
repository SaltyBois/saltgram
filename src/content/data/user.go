package data

import "github.com/go-playground/validator"

type User struct {
	Email          string    `json:"email" gorm:"primaryKey" validate:"required"`
	FullName       string    `json:"fullName" validate:"required"`
	Username       string    `json:"username" validate:"required"`
}

func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}