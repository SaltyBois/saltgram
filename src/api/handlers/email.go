package handlers

import (
	"saltgram/protos/email/premail"
	"saltgram/protos/users/prusers"

	"github.com/sirupsen/logrus"
)

type Email struct {
	l  *logrus.Logger
	ec premail.EmailClient
	uc prusers.UsersClient
}

func NewEmail(l *logrus.Logger, ec premail.EmailClient, uc prusers.UsersClient) *Email {
	return &Email{l: l, ec: ec, uc: uc}
}
