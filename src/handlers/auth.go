package handlers

import (
	"fmt"
	"log"

	"github.com/dgrijalva/jwt-go"
)

type Auth struct {
	l *log.Logger
}

type AccessClaims struct {
	Username       string             `json:"username"`
	Password       string             `json:"password"`
	StandardClaims jwt.StandardClaims `json:"standardClaims"`
}

type RefreshClaims struct {
	Username       string             `json:"username"`
	StandardClaims jwt.StandardClaims `json:"standardClaims"`
}

var ErrorEmptyClaims = fmt.Errorf("empty credentials")

func (uc AccessClaims) Valid() error {
	if len(uc.Username) <= 0 || len(uc.Password) <= 0 {
		return ErrorEmptyClaims
	}

	return uc.StandardClaims.Valid()
}

func (rc RefreshClaims) Valid() error {
	if len(rc.Username) <= 0 {
		return ErrorEmptyClaims
	}

	return rc.StandardClaims.Valid()
}

func NewAuth(l *log.Logger) *Auth {
	return &Auth{l}
}
