package data

import (
	"fmt"

	"github.com/go-playground/validator"
)

type Login struct {
	Username  string    `json:"username" validate:"required"`
	Password  string    `json:"password" validate:"required"`
	ReCaptcha ReCaptcha `json:"reCaptcha" validate:"required"`
}

// NOTE(Jovan): Refresh token
type Refresh struct {
	Username string `json:"username"`
	Token string `json:"token"`
}

func (l *Login) Validate() error {
	// TODO(Jovan): Extract into a global validator?
	validate := validator.New()
	return validate.Struct(l)
}

var refreshTokens = []*Refresh{}

func AddRefreshToken(username, token string) {
	refreshTokens = append(refreshTokens, &Refresh{username, token})
}

var ErrorRefreshTokenNotFound = fmt.Errorf("refresh token not found")

func GetRefreshToken(username string) (string, error) {
	for _, rt := range refreshTokens {
		if rt.Username == username {
			return rt.Token, nil
		}
	}
	return "", ErrorRefreshTokenNotFound
}

func GetRefreshTokens() []*Refresh {
	return refreshTokens
}

var ErrorInvalidRefreshToken = fmt.Errorf("invalid refresh token")

func (r *Refresh) Validate() error {
	for _, rt := range refreshTokens {
		if rt.Token == r.Token {
			return nil
		}
	}
	return ErrorInvalidRefreshToken
}
