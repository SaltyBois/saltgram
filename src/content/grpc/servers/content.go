package servers

import (
	"log"
	"saltgram/content/data"
	"saltgram/protos/content/prcontent"
)

type Content struct {
	prcontent.UnimplementedContentServer
	l  *log.Logger
	db *data.DBConn
}

func NewContent(l *log.Logger, db *data.DBConn) *Content {
	return &Content{
		l:  l,
		db: db,
	}
}
