package webutil_test

import (
	"net/http"
	"testing"

	"github.com/noypi/webutil"
	assertpkg "github.com/stretchr/testify/assert"
)

func TestMidSeq(t *testing.T) {
	assert := assertpkg.New(t)

	var callseq string
	h0 := func(nexth http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callseq += "h0->"
			nexth.ServeHTTP(w, r)
		})
	}
	h1 := func(a string, nexth http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callseq += "h1." + a + "->"
			nexth.ServeHTTP(w, r)
		})
	}
	h2 := func(a, b string, nexth http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callseq += "h2." + a + b + "->"
			nexth.ServeHTTP(w, r)
		})
	}

	h3 := func(a, b, c string, nexth http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callseq += "h3." + a + b + c + "->"
			nexth.ServeHTTP(w, r)
		})
	}
	h := func(w http.ResponseWriter, r *http.Request) {
		callseq += "h->"
	}

	fnH := webutil.MidSeq(
		http.HandlerFunc(h),
		webutil.MidFn(h2, "2a", "2b"),
		webutil.MidFn(h0),
		webutil.MidFn(h3, "3a", "3b", "3c"),
		webutil.MidFn(h1, "1a"),
	)

	fnH.ServeHTTP(nil, nil)
	assert.Equal("h2.2a2b->h0->h3.3a3b3c->h1.1a->h->", callseq)

}
