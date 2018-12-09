package log

import (
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

// Config represents the logging configuration
// that is set via environment variables.
type Config struct {
	LogLevel    string `envconfig:"LOG_LEVEL"`
	ServiceName string `envconfig:"SERVICE_NAME" default:"unknown"`
	Version     string `envconfig:"SERVICE_VERSION" default:"dev"`
}

// NewLogger returns a ready-to-use logger.
func NewLogger() *logrus.Entry {

	// get the service fields
	var conf Config
	if err := envconfig.Process("", &conf); err != nil {
		log.Fatal("could not load the logging environment variables: ", err)
	}

	// define the logger format and its properties
	logger := &logrus.Logger{
		Formatter: &logrus.TextFormatter{
			FullTimestamp:    true,
			QuoteEmptyFields: true,
			// RFC3339 with milliseconds
			TimestampFormat: "2006-01-02T15:04:05.999Z07:00",
		},
		Level: logrus.DebugLevel, // default
		Out:   os.Stdout,
	}

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	if conf.LogLevel != "" {
		lvl, err := logrus.ParseLevel(conf.LogLevel)
		if err != nil {
			logrus.WithError(err).Fatalf("could not set log level '%v'", conf.LogLevel)
		}
		logger.Level = lvl
	}

	// add service specific fields
	fields := logrus.Fields{
		"service": conf.ServiceName,
		"version": conf.Version,
	}

	log := logger.WithFields(fields)

	return log
}
