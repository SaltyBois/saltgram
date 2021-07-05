package handlers

import (
	"log"
)

type Emails struct {
	l *log.Logger
}

func NewEmail(l *log.Logger) *Emails {
	return &Emails{l}
}
