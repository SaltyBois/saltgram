package handlers

import (
	"saltgram/protos/notifications/prnotifications"
	"saltgram/protos/users/prusers"

	"github.com/sirupsen/logrus"
)

type Notification struct {
	l  *logrus.Logger
	nc prnotifications.NotificationsClient
	uc prusers.UsersClient
}

func NewNotification(l *logrus.Logger, nc prnotifications.NotificationsClient, uc prusers.UsersClient) *Notification {
	return &Notification{
		l:  l,
		nc: nc,
		uc: uc,
	}
}
