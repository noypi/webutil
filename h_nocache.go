package webutil

import (
	"net/http"
	"time"
)

func NoCache(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hdr := w.Header()
		hdr.Set("Cache-Control", "no-cache, private, max-age=0")
		hdr.Set("Expires", time.Unix(0, 0).Format(http.TimeFormat))
		hdr.Set("Pragma", "no-cache")
		hdr.Set("X-Accel-Expires", "0")

		h.ServeHTTP(w, r)
	})

}
