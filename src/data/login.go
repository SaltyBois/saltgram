package data

import "github.com/go-playground/validator"

type Login struct {
	Username  string    `json:"username" validate:"required"`
	Password  string    `json:"password" validate:"required"`
	ReCaptcha ReCaptcha `json:"reCaptcha" validate:"required"`
}

func (l *Login) Validate() error {
	// TODO(Jovan): Extract into a global validator?
	validate := validator.New()
	return validate.Struct(l)
}
