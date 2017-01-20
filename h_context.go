package webutil

import (
	"context"
	"net/http"
)

func WithCtxValue(key, value interface{}, nexth http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), key, value)
		nexth.ServeHTTP(w, r.WithContext(ctx))
	})
}
