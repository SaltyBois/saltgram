package handlers

import (
	"log"
	"saltgram/protos/auth/prauth"
)

type Auth struct {
	l          *log.Logger
	authClient prauth.AuthClient
}

func NewAuth(l *log.Logger, authClient prauth.AuthClient) *Auth {
	return &Auth{l: l, authClient: authClient}
}
