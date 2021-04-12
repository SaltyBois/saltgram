package handlers

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Login struct {
	l *log.Logger
}

type UserClaims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	NBF      string `json:"nbf"`
}

type KeyLogin struct{}

var ErrorEmptyClaims = fmt.Errorf("empty credentials")

func (uc *UserClaims) Valid() error {
	if len(uc.Username) <= 0 || len(uc.Password) <= 0 || len(uc.NBF) <= 0 {
		return ErrorEmptyClaims
	}
	
	return nil
}

var ErrorJWTExpired = fmt.Errorf("jwt expired")

func (uc *UserClaims) CheckDate() error {
	layout := os.Getenv("TIME_LAYOUT")
	t, err := time.Parse(layout, uc.NBF)
	if err != nil {
		return err
	}

	if t.UTC().Before(time.Now()) {
		return ErrorJWTExpired
	}
	return nil
}

func NewLogin(l *log.Logger) *Login {
	return &Login{l}
}
