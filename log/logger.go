package log

import (
	"github.com/devingen/api-core/dvnruntime"
	"github.com/sirupsen/logrus"
)

const (
	ContextKeyLogger dvnruntime.ContextKey = "context-key-logger"
)

var logger *logrus.Entry

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

func Warning(args ...interface{}) {
	logger.Warning(args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

func Panic(args ...interface{}) {
	logger.Panic(args...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}

func SetLevel(level logrus.Level) {
	logger.Level = level
}

func WithFields(fields logrus.Fields) *logrus.Entry {
	return logger.WithFields(fields)
}

var formatter = logrus.JSONFormatter{
	FieldMap: logrus.FieldMap{
		logrus.FieldKeyMsg: "message",
	},
}

func Init() {
	logBase := logrus.New()
	logBase.SetFormatter(&formatter)
	logger = logBase.WithFields(logrus.Fields{})
}

func InitWithBaseFields(baseFields logrus.Fields) {
	logBase := logrus.New()
	logBase.SetFormatter(&formatter)
	logger = logBase.WithFields(baseFields)
}
