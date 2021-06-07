package handlers

import (
	"log"
)

type Content struct {
	l *log.Logger
}

func NewContent(l *log.Logger) *Content {
	return &Content{l}
}
