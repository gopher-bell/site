package main

import (
	"net/http"

	"github.com/gopher-bell/site/log"
	"go.uber.org/zap"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")

		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.ZapLogger.Info("requesting to...", zap.String("remote addr", r.RemoteAddr), zap.String("proto", r.Proto), zap.String("method", r.Method), zap.String("url", r.URL.RequestURI()))
		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				log.ZapLogger.Error("panic occured", zap.Any("err", err))
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}()

		next.ServeHTTP(w, r)
	})
}
