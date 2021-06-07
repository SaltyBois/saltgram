package servers

import
(
	"log"
	"saltgram/content/data"
)

type Content struct {
	l  *log.Logger
	db *data.DBConn
}

func NewContent(l *log.Logger, db *data.DBConn) *Content {
	return &Content{
		l:  l,
		db: db,
	}
}
