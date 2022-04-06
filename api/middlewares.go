package api

import (
	"context"
	"github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"
	application "github.com/ukrainian-brothers/board-backend/app"
	"github.com/ukrainian-brothers/board-backend/internal/common"
	"net/http"
	"runtime/debug"
	"time"
)

type MiddlewareProvider struct {
	sessionStore *sessions.CookieStore
	app          *application.Application
	cfg          *common.Config
}

func NewMiddlewareProvider(sessionStore *sessions.CookieStore, app *application.Application, cfg *common.Config) *MiddlewareProvider {
	return &MiddlewareProvider{sessionStore: sessionStore, app: app, cfg: cfg}
}

func (p MiddlewareProvider) AuthMiddleware(next http.HandlerFunc, logger *log.Entry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := p.sessionStore.Get(r, p.cfg.Session.SessionKey)
		if err != nil {
			logger.WithError(err).Error("failed getting session from store")
			next.ServeHTTP(w, r)
			return
		}

		if session.Values["user_login"] != nil {
			// Read session login and put it into ctx - so later on it can be used to verify if user has access to an entity
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "user_login", session.Values["user_login"].(string))))
			return
		}
		
		next.ServeHTTP(w, r)
	}
}

func (p MiddlewareProvider) LoggingMiddleware(logger *log.Entry) func(http.Handler) http.Handler {
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

func (p MiddlewareProvider) BodyLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" || r.Method == "PUT" {
			r.Body = http.MaxBytesReader(w, r.Body, 1048576)
		}
		next.ServeHTTP(w, r)
	})
}
