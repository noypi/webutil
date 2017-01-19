package webutil

import (
	"net/http"
)

type MidInfo struct {
	fn     interface{}
	params []string
}

func MidFn(fn interface{}, params ...string) *MidInfo {
	return &MidInfo{
		fn:     fn,
		params: params,
	}
}

func MidSeqFunc(handlerToWrap http.HandlerFunc, fs ...*MidInfo) http.Handler {
	return MidSeq(handlerToWrap, fs...)
}

func MidSeq(handlerToWrap http.Handler, fs ...*MidInfo) http.Handler {
	var currfn http.Handler = handlerToWrap
	for i := len(fs) - 1; i >= 0; i-- {
		var f *MidInfo = fs[i]
		switch fn := f.fn.(type) {
		case fn0:
			currfn = fn(currfn)
		case fn1:
			currfn = fn(f.params[0], currfn)
		case fn2:
			currfn = fn(f.params[0], f.params[1], currfn)
		case fn3:
			currfn = fn(f.params[0], f.params[1], f.params[2], currfn)
		case fn4:
			currfn = fn(f.params[0], f.params[1], f.params[2], f.params[3], currfn)
		case fn5:
			currfn = fn(f.params[0], f.params[1], f.params[2], f.params[3], f.params[4], currfn)
		case fn6:
			currfn = fn(f.params[0], f.params[1], f.params[2], f.params[3], f.params[4], f.params[5], currfn)
		}
	}

	return currfn
}

type fn0 func(nexth http.Handler) http.Handler
type fn1 func(a string, nexth http.Handler) http.Handler
type fn2 func(a, b string, nexth http.Handler) http.Handler
type fn3 func(a, b, c string, nexth http.Handler) http.Handler
type fn4 func(a, b, c, d string, nexth http.Handler) http.Handler
type fn5 func(a, b, c, d, e string, nexth http.Handler) http.Handler
type fn6 func(a, b, c, d, e, f string, nexth http.Handler) http.Handler
