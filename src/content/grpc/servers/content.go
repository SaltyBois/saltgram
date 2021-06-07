package servers

import "log"

type Content struct {
	l  *log.Logger
	db *data.DBConn
}

func NewUsers(l *log.Logger, db *data.DBConn) *Content {
	return &Content{
		l:  l,
		db: db,
	}
}
