package handlers

import (
	"saltgram/protos/admin/pradmin"
	"saltgram/protos/users/prusers"

	"github.com/sirupsen/logrus"
)

type Admin struct {
	l  *logrus.Logger
	ac pradmin.AdminClient
	uc prusers.UsersClient
}

func NewAdmin(l *logrus.Logger, ac pradmin.AdminClient, uc prusers.UsersClient) *Admin {
	return &Admin{l: l, ac: ac, uc: uc}
}
