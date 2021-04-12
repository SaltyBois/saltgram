package handlers

import (
	"log"
)

type Login struct {
	l *log.Logger
}

type KeyLogin struct{}

func NewLogin(l *log.Logger) *Login {
	return &Login{l}
}
