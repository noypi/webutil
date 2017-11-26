package webutil

import (
	"net/http"

	"github.com/noypi/router"
)

func IfSuccessRedirect(theURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := router.ContextR(r)
		if _, hasError := c.Get(ErrorKey); !hasError {
			c.Redirect(http.StatusTemporaryRedirect, theURL)
		}
	}
}
