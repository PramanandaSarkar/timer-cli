package utils

import (
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func InitLogger(logLevel string) {
	Logger = logrus.New()
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		Logger.SetLevel(logrus.InfoLevel)
	} else {
		Logger.SetLevel(level)
	}
	Logger.SetFormatter(&logrus.JSONFormatter{})
}