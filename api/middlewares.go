package api

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"runtime/debug"
	"time"
)

// responseWriter is a minimal wrapper for http.ResponseWriter that allows the
// written HTTP status code to be captured for logging.
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true

	return
}

// LoggingMiddleware logs the incoming HTTP request & its duration.
func LoggingMiddleware(logger *log.Entry) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					logger.WithFields(log.Fields{
						"err":   err,
						"trace": debug.Stack(),
					}).Info("unknown internal error")
				}
			}()

			start := time.Now()
			wrapped := wrapResponseWriter(w)
			next.ServeHTTP(wrapped, r)
			logger.WithFields(log.Fields{
				"status":   wrapped.status,
				"method":   r.Method,
				"path":     r.URL.EscapedPath(),
				"duration": time.Since(start),
			}).Info()
		}

		return http.HandlerFunc(fn)
	}
}

func BodyLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" || r.Method == "PUT" {
			r.Body = http.MaxBytesReader(w, r.Body, 1048576)
		}
		next.ServeHTTP(w, r)
	})
}
