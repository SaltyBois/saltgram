package s3

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/sirupsen/logrus"
)

type S3 struct {
	l *logrus.Logger
	s *session.Session
}

