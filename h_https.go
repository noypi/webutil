package webutil

import (
	"net/http"
	"strings"

	"github.com/noypi/router"
)

func EnsureHTTPS(redirect string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.Proto, "HTTPS") {
			c := router.ContextW(w)
			http.Redirect(w, r, redirect, http.StatusTemporaryRedirect)
			c.Abort()
			return
		}
	}
}
