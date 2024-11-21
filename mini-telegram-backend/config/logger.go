package config

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func InitializeLogger() {
	Logger = logrus.New()

	Logger.SetFormatter(&logrus.JSONFormatter{})

	logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logrus.Fatalf("Failed to open log file: %v", err)
	}

	multiWriter := io.MultiWriter(logFile, os.Stdout)
	Logger.SetOutput(multiWriter)

	// Set log level
	level, err := logrus.ParseLevel(GetEnv("LOG_LEVEL", "info"))
	if err != nil {
		level = logrus.InfoLevel
	}
	Logger.SetLevel(level)

	Logger.Info("Logger initialized and writing to both app.log and console")
}
