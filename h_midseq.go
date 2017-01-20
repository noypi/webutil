package webutil

import (
	"net/http"
)

type MidInfo struct {
	Fn     interface{}
	Params []interface{}
}

func MidFn(fn interface{}, params ...string) *MidInfo {
	paramsif := make([]interface{}, len(params))
	for i, v := range params {
		paramsif[i] = v
	}
	return &MidInfo{
		Fn:     fn,
		Params: paramsif,
	}
}

func MidFnIf(fn interface{}, params ...interface{}) *MidInfo {
	return &MidInfo{
		Fn:     fn,
		Params: params,
	}
}

func MidLogFn(fn interface{}, a LogFunc) *MidInfo {
	return MidFnIf(fn, a)
}

func MidSeqFunc(handlerToWrap http.HandlerFunc, fs ...*MidInfo) http.Handler {
	return MidSeq(handlerToWrap, fs...)
}

var DefaultMidSeq func(fn interface{}, params []interface{}) http.Handler

func MidSeq(handlerToWrap http.Handler, fs ...*MidInfo) http.Handler {
	var currfn http.Handler = handlerToWrap
	for i := len(fs) - 1; i >= 0; i-- {
		var f *MidInfo = fs[i]
		var ps = fs[i].Params
		switch fn := f.Fn.(type) {
		// interface{} params
		case fn1If:
			currfn = fn(ps[0], currfn)
		case fn2If:
			currfn = fn(ps[0], ps[1], currfn)
		case fn3If:
			currfn = fn(ps[0], ps[1], ps[2], currfn)

		// string params
		case fn0:
			currfn = fn(currfn)
		case fn1:
			currfn = fn(ps[0].(string), currfn)
		case fn2:
			currfn = fn(ps[0].(string), ps[1].(string), currfn)
		case fn3:
			currfn = fn(ps[0].(string), ps[1].(string), ps[2].(string), currfn)
		case fn4:
			currfn = fn(ps[0].(string), ps[1].(string), ps[2].(string), ps[3].(string), currfn)
		case fn5:
			currfn = fn(ps[0].(string), ps[1].(string), ps[2].(string), ps[3].(string), ps[4].(string), currfn)
		case fn6:
			currfn = fn(ps[0].(string), ps[1].(string), ps[2].(string), ps[3].(string), ps[4].(string), ps[5].(string), currfn)

		// others
		case func(LogFunc, http.Handler) http.HandlerFunc:
			currfn = fn(ps[0].(LogFunc), currfn)
		default:
			if nil != DefaultMidSeq {
				currfn = DefaultMidSeq(fn, ps)
			}
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

type fn1If func(a interface{}, nexth http.Handler) http.Handler
type fn2If func(a, b interface{}, nexth http.Handler) http.Handler
type fn3If func(a, b, c interface{}, nexth http.Handler) http.Handler
