package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type loggerKey struct{}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sl := slog.New(slog.NewTextHandler(log.Logger, nil))

		ctx := context.WithValue(r.Context(), loggerKey{}, sl)
		r = r.WithContext(ctx)

		start := time.Now()
		rw := &responseWriter{ResponseWriter: w}
		next.ServeHTTP(rw, r)
		latency := time.Since(start)

		log.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", rw.status).
			Str("client_ip", r.RemoteAddr).
			Dur("latency", latency).
			Int("bytes", rw.size).
			Send()
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.size += n
	if rw.status == 0 {
		rw.status = http.StatusOK
	}
	return n, err
}

func GetLogger(r *http.Request) *slog.Logger {
	l, _ := r.Context().Value(loggerKey{}).(*slog.Logger)
	return l
}
