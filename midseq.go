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

func MidSeq(handlerToWrap http.HandlerFunc, fs ...*MidInfo) http.HandlerFunc {
	var currfn http.HandlerFunc = handlerToWrap
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

type fn0 func(nexth http.HandlerFunc) http.HandlerFunc
type fn1 func(a string, nexth http.HandlerFunc) http.HandlerFunc
type fn2 func(a, b string, nexth http.HandlerFunc) http.HandlerFunc
type fn3 func(a, b, c string, nexth http.HandlerFunc) http.HandlerFunc
type fn4 func(a, b, c, d string, nexth http.HandlerFunc) http.HandlerFunc
type fn5 func(a, b, c, d, e string, nexth http.HandlerFunc) http.HandlerFunc
type fn6 func(a, b, c, d, e, f string, nexth http.HandlerFunc) http.HandlerFunc
