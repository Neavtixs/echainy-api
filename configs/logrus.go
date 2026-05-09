package configs

import (
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

func NewLogger() *logrus.Logger {
	log := logrus.New()

	log.SetOutput(os.Stdout)
	log.SetLevel(parseLogLevel(os.Getenv("LOG_LEVEL")))
	log.SetFormatter(newLogFormatter(os.Getenv("BE_ENV")))

	return log
}

func newLogFormatter(env string) logrus.Formatter {
	if env == "production" {
		return &logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		}
	}

	return &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "15:04:05",
		ForceColors:     true,
	}
}

func parseLogLevel(value string) logrus.Level {
	if value == "" {
		return logrus.InfoLevel
	}

	level, err := logrus.ParseLevel(value)
	if err == nil {
		return level
	}

	number, err := strconv.Atoi(value)
	if err != nil {
		return logrus.InfoLevel
	}

	return logrus.Level(number)
}
