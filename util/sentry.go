package util

import (
	"net/http"

	"github.com/getsentry/sentry-go"

	"github.com/gorilla/context"
)

//
// SentryWrapper wraps http handlers with sentry sdk
// to log errors
//
func SentryWrapper(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := context.Get(r, "token")
		defer func() {
			if recover := recover(); r != nil {
				sentry.ConfigureScope(func(scope *sentry.Scope) {
					scope.SetExtra("user_token", token)
				})
				if recover, ok := recover.(error); ok {
					ErrorResponder(w, http.StatusInternalServerError, recover)
				}
			}
		}()
		next(w, r)
	}
}
