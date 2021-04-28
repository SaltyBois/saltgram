package handlers

import (
	"log"
	"saltgram/protos/email/premail"
	"saltgram/protos/users/prusers"
)

type Email struct {
	l  *log.Logger
	ec premail.EmailClient
	uc prusers.UsersClient
}

func NewEmail(l *log.Logger, ec premail.EmailClient, uc prusers.UsersClient) *Email {
	return &Email{l: l, ec: ec, uc: uc}
}
