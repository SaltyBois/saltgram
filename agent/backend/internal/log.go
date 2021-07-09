package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	L *logrus.Logger
	f *os.File
}

func NewLogger(serviceName string) *Logger {
	logger := &Logger{}

	logger.L = logrus.New()
	logger.L.SetReportCaller(true)
	logger.L.SetOutput(os.Stdout)
	logger.L.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "01-02-2006 15:04:05",
		FullTimestamp:   true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return f.Function, fmt.Sprintf("%s:%d", formatFilePath(f.File), f.Line)
		}})
	logger.L.SetLevel(logrus.InfoLevel)
	return logger
}

func formatFilePath(path string) string {
	arr := strings.Split(filepath.ToSlash(path), "/")
	return arr[len(arr)-1]
}

func (l *Logger) Close() {
	l.f.Close()
}
