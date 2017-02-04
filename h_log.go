package webutil

import (
	"net/http"

	"github.com/noypi/util"
)

func WithErrLogger(fn util.LogFunc, nexth http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := util.WithErrLogger(r.Context(), fn)
		nexth.ServeHTTP(w, r.WithContext(ctx))
	})
}

func WithInfoLogger(fn func(fmt string, params ...interface{}), nexth http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := util.WithInfoLogger(r.Context(), fn)
		nexth.ServeHTTP(w, r.WithContext(ctx))
	})
}
