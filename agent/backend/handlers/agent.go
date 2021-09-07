package handlers

import (
	"agent/data"

	"github.com/sirupsen/logrus"
)

type Agent struct {
	l  *logrus.Logger
	db *data.DBConn
}

func NewAgent(l *logrus.Logger, db *data.DBConn) *Agent {
	return &Agent{l: l, db: db}
}
