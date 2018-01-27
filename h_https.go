package webutil

import (
	"net/http"

	"github.com/noypi/router"
)

func EnsureHTTPS(w http.ResponseWriter, r *http.Request) {
	if "127.0.0.1" == r.RemoteAddr {
		return
	}
	if "true" == r.Header.Get("X-Appengine-Cron") ||
		0 < len(r.Header.Get("X-AppEngine-QueueName")) {
		return
	}
	if r.URL.Scheme != "https" {
		c := router.ContextW(w)
		urlpath := "https://" + r.URL.Host + r.URL.Path
		http.Redirect(w, r, urlpath, http.StatusTemporaryRedirect)
		c.Abort()
		return
	}
}
