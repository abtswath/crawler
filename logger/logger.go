package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

func New() *logrus.Logger {
	log := logrus.New()
	log.SetLevel(logrus.TraceLevel)
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.TextFormatter{})
	return log
}
