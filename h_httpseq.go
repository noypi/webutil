package webutil

import (
	"net/http"

	"github.com/noypi/router"
)

func HttpSequence(finally http.HandlerFunc, h ...http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := router.ContextR(r)
		for i := 0; i < len(h) && !c.IsAborted(); i++ {
			h[i].ServeHTTP(c.Writer, c.Request)
		}
		finally.ServeHTTP(c.Writer, c.Request)
	}
}