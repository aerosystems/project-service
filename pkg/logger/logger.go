package logger

import (
	"fmt"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

type Logger struct{ *logrus.Logger }

func (l *Logger) Say(msg string) {
	l.Info(msg)
}
func (l *Logger) Sayf(fmt string, args ...interface{}) {
	l.Infof(fmt, args)
}
func (l *Logger) SayWithField(msg string, k string, v interface{}) {
	l.WithField(k, v).Info(msg)
}
func (l *Logger) SayWithFields(msg string, fields map[string]interface{}) {
	l.WithFields(fields).Info(msg)
}

func NewLogger() *Logger {
	filename := os.Getenv("HOSTNAME")
	logLevel := logrus.InfoLevel
	log := logrus.New()
	log.SetLevel(logLevel)
	rotateFileHook, err := NewRotateFileHook(RotateFileConfig{
		Filename:   fmt.Sprintf("/app/logs/%s.log", filename),
		MaxSize:    50, // megabytes
		MaxBackups: 3,  // amounts
		MaxAge:     28, //days
		Level:      logLevel,
		Formatter: &logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		},
	})

	if err != nil {
		logrus.Fatalf("Failed to initialize file rotate hook: %v", err)
	}

	log.SetOutput(colorable.NewColorableStdout())
	log.SetFormatter(&logrus.TextFormatter{
		PadLevelText:    true,
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05.000Z",
	})

	log.AddHook(rotateFileHook)

	return &Logger{log}
}
