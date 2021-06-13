package servers

import (
	"context"
	"saltgram/notification/data"
	"saltgram/protos/notification/prnotification"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Notification struct {
	prnotification.UnimplementedNotificationServer
	l  *logrus.Logger
	db *data.DBConn
}

func NewNotification(l *logrus.Logger, db *data.DBConn) *Notification {
	return &Notification{
		l:  l,
		db: db,
	}
}