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
		switch len(fs[i].params) {
		case 0:
			currfn = f.fn.(fn0)(currfn)
		case 1:
			currfn = f.fn.(fn1)(f.params[0], currfn)
		case 2:
			currfn = f.fn.(fn2)(f.params[0], f.params[1], currfn)
		case 3:
			currfn = f.fn.(fn3)(f.params[0], f.params[1], f.params[2], currfn)
		case 4:
			currfn = f.fn.(fn4)(f.params[0], f.params[1], f.params[2], f.params[3], currfn)
		case 5:
			currfn = f.fn.(fn5)(f.params[0], f.params[1], f.params[2], f.params[3], f.params[4], currfn)
		case 6:
			currfn = f.fn.(fn6)(f.params[0], f.params[1], f.params[2], f.params[3], f.params[4], f.params[5], currfn)
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
