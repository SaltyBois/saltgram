package handlers

import (
	"log"
	"saltgram/protos/users/prusers"
)

type Users struct {
	l  *log.Logger
	uc prusers.UsersClient
}

func NewUsers(l *log.Logger, uc prusers.UsersClient) *Users {
	return &Users{l: l, uc: uc}
}
