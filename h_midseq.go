package webutil

import (
	"log"
	"net/http"
	"reflect"
	"runtime"

	"github.com/noypi/logfn"
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

func MidLogFn(fn interface{}, a logfn.LogFunc) *MidInfo {
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
		case func(a interface{}, nexth http.Handler) http.Handler:
			currfn = fn(ps[0], currfn)
		case func(a, b interface{}, nexth http.Handler) http.Handler:
			currfn = fn(ps[0], ps[1], currfn)
		case func(a, b, c interface{}, nexth http.Handler) http.Handler:
			currfn = fn(ps[0], ps[1], ps[2], currfn)

		// string params
		case func(http.Handler) http.Handler:
			currfn = fn(currfn)
		case func(a string, nexth http.Handler) http.Handler:
			currfn = fn(ps[0].(string), currfn)
		case func(a, b string, nexth http.Handler) http.Handler:
			currfn = fn(ps[0].(string), ps[1].(string), currfn)
		case func(a, b, c string, nexth http.Handler) http.Handler:
			currfn = fn(ps[0].(string), ps[1].(string), ps[2].(string), currfn)
		case func(a, b, c, d string, nexth http.Handler) http.Handler:
			currfn = fn(ps[0].(string), ps[1].(string), ps[2].(string), ps[3].(string), currfn)
		case func(a, b, c, d, e string, nexth http.Handler) http.Handler:
			currfn = fn(ps[0].(string), ps[1].(string), ps[2].(string), ps[3].(string), ps[4].(string), currfn)
		case func(a, b, c, d, e, f string, nexth http.Handler) http.Handler:
			currfn = fn(ps[0].(string), ps[1].(string), ps[2].(string), ps[3].(string), ps[4].(string), ps[5].(string), currfn)

		// others
		case func(logfn.LogFunc, http.Handler) http.HandlerFunc:
			currfn = fn(ps[0].(logfn.LogFunc), currfn)
		default:
			log.Println("default: ", runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name(), ", params len=", len(ps))
			if nil != DefaultMidSeq {
				currfn = DefaultMidSeq(fn, ps)
			}
		}
	}

	return currfn
}
