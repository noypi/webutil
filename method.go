package webutil

import (
	"net/http"
)

func Method(method string, nexth http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			return
		}
		nexth.ServeHTTP(w, r)
	})
}
