package handlers

import (
	"log"
	"saltgram/protos/auth/prauth"
)

type Auth struct {
	l  *log.Logger
	ac prauth.AuthClient
}

func NewAuth(l *log.Logger, ac prauth.AuthClient) *Auth {
	return &Auth{l: l, ac: ac}
}
