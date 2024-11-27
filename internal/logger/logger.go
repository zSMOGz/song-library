package logger

import (
	"os"

	"github.com/rs/zerolog"
)

func NewLogger() zerolog.Logger {
	return zerolog.New(os.Stdout).
		With().
		Timestamp().
		Caller().
		Logger()
}
