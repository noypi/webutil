package webutil

import (
	"net/http"
)

func RedirectHandler(urlstr string, code int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, urlstr, code)
	})
}
