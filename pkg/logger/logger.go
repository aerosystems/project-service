package logger

import (
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
)

type Logger struct{ *logrus.Logger }

func NewLogger() *Logger {
	logLevel := logrus.InfoLevel
	log := logrus.New()
	log.SetLevel(logLevel)
	log.SetOutput(colorable.NewColorableStdout())
	log.SetFormatter(&logrus.TextFormatter{
		PadLevelText:    true,
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05.000Z",
	})
	return &Logger{log}
}
