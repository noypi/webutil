package webutil

import (
	"context"
	"log"
	"net/http"
)

type LogFunc func(fmt string, params ...interface{})

const (
	LogErrName  = "$errlog"
	LogInfoName = "$infolog"
)

func LogErr(ctx context.Context, fmt string, params ...interface{}) {
	var fn LogFunc = log.Printf
	if nil != ctx.Value(LogErrName) {
		fn = ctx.Value(LogErrName).(LogFunc)
	}

	fn(fmt, params...)
}

func LogInfo(ctx context.Context, fmt string, params ...interface{}) {
	var fn LogFunc = log.Printf
	if nil != ctx.Value(LogInfoName) {
		fn = ctx.Value(LogInfoName).(LogFunc)
	}

	fn(fmt, params...)
}

func WithErrLogger(fn LogFunc, nexth http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), LogErrName, fn)
		nexth.ServeHTTP(w, r.WithContext(ctx))
	})
}

func WithInfoLogger(fn LogFunc, nexth http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), LogInfoName, fn)
		nexth.ServeHTTP(w, r.WithContext(ctx))
	})
}
