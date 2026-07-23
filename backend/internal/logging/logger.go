package logging

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

func NewLogger(logFormat, env string, writers ...io.Writer) zerolog.Logger {
	var output io.Writer = os.Stdout
	if len(writers) > 0 {
		output = writers[0]
	}

	format := logFormat
	if format == "auto" {
		if env == "production" {
			format = "json"
		} else {
			format = "text"
		}
	}

	if format == "text" {
		output = zerolog.ConsoleWriter{
			Out:        output,
			TimeFormat: time.RFC3339,
		}
	}

	return zerolog.New(output).With().Timestamp().Logger()
}
