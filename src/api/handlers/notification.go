package handlers

import (
	"saltgram/protos/notification/prnotification"
	"saltgram/protos/users/prusers"

	"github.com/sirupsen/logrus"
)

type Notification struct {
	l  *logrus.Logger
	nc prnotification.NotificationClient
	uc prusers.UsersClient
}

func NewNotification(l *logrus.Logger, nc prnotification.NotificationClient, uc prusers.UsersClient) *Notification {
	return &Notification{l: l, nc: nc, uc: uc}
}