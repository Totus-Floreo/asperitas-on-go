package middleware

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func AccessLog(logger *zap.SugaredLogger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			logger.Infow("New Request",
				"method", r.Method,
				"remote_addr", r.RemoteAddr,
				"body", r.Body,
				"url", r.URL.Path,
				"time", time.Since(start),
			)
		})
	}
}
