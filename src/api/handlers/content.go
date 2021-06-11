package handlers

import (
	"saltgram/protos/content/prcontent"
	"saltgram/protos/users/prusers"

	"github.com/sirupsen/logrus"
)

type Content struct {
	l  *logrus.Logger
	cc prcontent.ContentClient
	uc prusers.UsersClient
}

func NewContent(l *logrus.Logger, cc prcontent.ContentClient, uc prusers.UsersClient) *Content {
	return &Content{l: l, cc: cc, uc: uc}
}
