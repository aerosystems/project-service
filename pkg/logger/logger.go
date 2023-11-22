package logger

import (
	"fmt"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"time"
)

type logger struct{ *logrus.Logger }

func (l *logger) Say(msg string) {
	l.Info(msg)
}
func (l *logger) Sayf(fmt string, args ...interface{}) {
	l.Infof(fmt, args)
}
func (l *logger) SayWithField(msg string, k string, v interface{}) {
	l.WithField(k, v).Info(msg)
}
func (l *logger) SayWithFields(msg string, fields map[string]interface{}) {
	l.WithFields(fields).Info(msg)
}

var l *logger

func NewLogger(filename string) *logger {
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

	// setting the global logger instance as singleton pattern
	l := &logger{log}
	_ = l

	return &logger{log}
}

// GetLogger returns the logger instance by singleton pattern
func GetLogger() *logger {
	return l
}
