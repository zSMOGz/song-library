package middleware

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"

	"song-library/internal/constants"
)

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func RequestLogger(log zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			sw := &statusWriter{ResponseWriter: w}
			next.ServeHTTP(sw, r)

			log.Info().
				Str(constants.LogFieldMethod, r.Method).
				Str(constants.LogFieldPath, r.URL.Path).
				Int(constants.LogFieldStatus, sw.status).
				Dur(constants.LogFieldDuration, time.Since(start)).
				Msg(constants.LogMsgRequest)
		})
	}
}
