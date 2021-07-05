package handlers

import (
	"saltgram/protos/auth/prauth"

	"github.com/sirupsen/logrus"
)

type Auth struct {
	l  *logrus.Logger
	ac prauth.AuthClient
}

func NewAuth(l *logrus.Logger, ac prauth.AuthClient) *Auth {
	return &Auth{l: l, ac: ac}
}
