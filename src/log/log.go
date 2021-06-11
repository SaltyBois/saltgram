package log

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
)

// NOTE(Jovan): Would make life easier in the long run to use
// custom Print, Warn, Error, etc. functions, but more hassle
// at startup

type Logger struct {
	L *logrus.Logger
	f *os.File
}

func NewLogger(serviceName string) *Logger {
	logger := &Logger {}
	err := makeDirectoryIfNotExists(filepath.FromSlash("../log/logs/" + serviceName))
	if err != nil {
		logrus.Fatalf("Ironically logrus failed to create dir: %v\n", err)
	}
	filename := filepath.FromSlash("../log/logs/" + serviceName + "/" + serviceName + ".log")
	rotatingHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename: filename,
		MaxSize: 50,
		MaxBackups: 5,
		Level: logrus.InfoLevel,
		Formatter: &logrus.JSONFormatter{
			TimestampFormat: "01-02-2006 15:04:05",
			DataKey: "data",
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				return f.Function, fmt.Sprintf("%s:%d", formatFilePath(f.File), f.Line)
		}},
	})
	if err != nil {
		logrus.Fatalf("Ironically logrus failed to create the rotating hook: %v\n", err)
	}
	logger.L = logrus.New()
	logger.L.AddHook(rotatingHook)
	logger.L.SetReportCaller(true)
	logger.L.SetOutput(os.Stdout)
	logger.L.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "01-02-2006 15:04:05",
		FullTimestamp: true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return f.Function, fmt.Sprintf("%s:%d", formatFilePath(f.File), f.Line)
		}})
	logger.L.SetLevel(logrus.InfoLevel)
	return logger
}

func makeDirectoryIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, os.ModeDir|0755)
	}
	return nil
}

func formatFilePath(path string) string {
	arr := strings.Split(filepath.ToSlash(path), "/")
	return arr[len(arr) - 1]
}

func (l *Logger) Close()  {
	l.f.Close()
}
