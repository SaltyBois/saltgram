package handlers

import (
	"log"
	"saltgram/protos/content/prcontent"
	"saltgram/protos/users/prusers"
)

type Content struct {
	l  *log.Logger
	cc prcontent.ContentClient
	uc prusers.UsersClient
}

func NewContent(l *log.Logger, cc prcontent.ContentClient, uc prusers.UsersClient) *Content {
	return &Content{l: l, cc: cc, uc: uc}
}
