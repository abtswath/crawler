package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

func New(level logrus.Level) *logrus.Logger {
	log := logrus.New()
	log.SetLevel(level)
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.TextFormatter{})
	return log
}
