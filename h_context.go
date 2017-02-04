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

func WithCtxKVMap(m map[interface{}]interface{}, nexth http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ctx = r.Context()
		for k, v := range m {
			ctx = context.WithValue(ctx, k, v)
		}
		nexth.ServeHTTP(w, r.WithContext(ctx))
	})
}
