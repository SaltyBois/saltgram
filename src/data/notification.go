package data

import "github.com/go-playground/validator"

type Notification struct {
	Username                  string `json:"username" validate:"required"`
	ReferredUsername          string `json:"referredUsername" validate:"required"`
	ReferredProfilePictureURL string `json:"profilePictureURL"`
	Type                      string `json:"type" validate:"required"`
	Seen                      bool   `json:"seen" validate:"required"`
}

func (n *Notification) Validate() error {
	// TODO(Jovan): Extract into a global validator?
	validate := validator.New()
	return validate.Struct(n)
}
