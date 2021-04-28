package data

import "github.com/go-playground/validator"

type UserDTO struct {
	Email     string    `json:"email" validate:"required"`
	FullName  string    `json:"fullName" validate:"required"`
	Username  string    `json:"username" validate:"required"`
	Password  string    `json:"password" validate:"required"`
	ReCaptcha ReCaptcha `json:"reCaptcha" validate:"required"`
}

func (u *UserDTO) Validate() error {
	// TODO(Jovan): Extract into a global validator?
	validate := validator.New()
	return validate.Struct(u)
}
