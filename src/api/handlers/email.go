package handlers

import (
	"log"
	"saltgram/protos/email/premail"
)

type Email struct {
	l  *log.Logger
	ec premail.EmailClient
}

func NewEmail(l *log.Logger, ec premail.EmailClient) *Email {
	return &Email{l: l, ec: ec}
}
