package data

import "github.com/go-playground/validator"

type UserDTO struct {
	Email     string    `json:"email" validate:"required"`
	FullName  string    `json:"fullName" validate:"required"`
	Username  string    `json:"username" validate:"required"`
	Password  string    `json:"password" validate:"required"`
	ReCaptcha ReCaptcha `json:"reCaptcha" validate:"required"`
}

type ProflieDTO struct {
	Email    string `json:"email" validate:"required"`
	FullName string `json:"fullName" validate:"required"`
	Username string `json:"username" validate:"required"`
	Public   bool   `json:"public" validate:"required"`
	Taggable bool   `json:"taggable" validate:"required"`
}

type FollowDTO struct {
	ProfileToFollow string `json:"profile" validate:"required"`
}

func (p *ProflieDTO) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

func (u *UserDTO) Validate() error {
	// TODO(Jovan): Extract into a global validator?
	validate := validator.New()
	return validate.Struct(u)
}

func (f *FollowDTO) Validate() error {
	validate := validator.New()
	return validate.Struct(f)
}
