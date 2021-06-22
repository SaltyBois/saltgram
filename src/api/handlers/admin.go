package handlers

import (
	"saltgram/protos/admin/pradmin"
	"saltgram/protos/content/prcontent"
	"saltgram/protos/users/prusers"

	"github.com/sirupsen/logrus"
)

type Admin struct {
	l  *logrus.Logger
	ac pradmin.AdminClient
	uc prusers.UsersClient
	cc prcontent.ContentClient
}

func NewAdmin(l *logrus.Logger, ac pradmin.AdminClient, uc prusers.UsersClient, cc prcontent.ContentClient) *Admin {
	return &Admin{l: l, ac: ac, uc: uc, cc: cc}
}
