package webutil

import (
	"net/http"

	"context"

	"github.com/noypi/logfn"
	"github.com/noypi/router"
)

func WithErrLogger(ctx context.Context, fn logfn.LogFunc) context.Context {
	c := ToStore(ctx)
	c.Set(router.LogErrName, fn)
	return ctx
}

func WithWarnLogger(ctx context.Context, fn logfn.LogFunc) context.Context {
	c := ToStore(ctx)
	c.Set(router.LogWarnName, fn)
	return ctx
}

func WithInfoLogger(ctx context.Context, fn logfn.LogFunc) context.Context {
	c := ToStore(ctx)
	c.Set(router.LogInfoName, fn)
	return ctx
}

func WithDebugLogger(ctx context.Context, fn logfn.LogFunc) context.Context {
	c := ToStore(ctx)
	c.Set(router.LogDebugName, fn)
	return ctx
}

func LogErr(ctx context.Context, fmt string, params ...interface{}) {
	router.GetErrLog(ctx)(fmt, params...)
}

func LogInfo(ctx context.Context, fmt string, params ...interface{}) {
	router.GetInfoLog(ctx)(fmt, params...)
}

func LogWarn(ctx context.Context, fmt string, params ...interface{}) {
	router.GetWarnLog(ctx)(fmt, params...)
}

func LogDebug(ctx context.Context, fmt string, params ...interface{}) {
	router.GetDebugLog(ctx)(fmt, params...)
}

func AddLoggerHandler(fnInfo, fnErr, fnWarn logfn.LogFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := router.ContextW(w)
		if nil != fnInfo {
			WithInfoLogger(ctx, fnInfo)
		} else if nil != fnErr {
			WithErrLogger(ctx, fnErr)
		} else if nil != fnWarn {
			WithWarnLogger(ctx, fnWarn)
		}
	})
}

func AddDebugLoggerHandler(fnDebug logfn.LogFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := router.ContextW(w)
		WithWarnLogger(ctx, fnDebug)
	})
}

func ErrLoggerHandler(fn logfn.LogFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := router.ContextW(w)
		WithErrLogger(ctx, fn)
	})
}

func InfoLoggerHandler(fn logfn.LogFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WithInfoLogger(router.ContextW(w), fn)
	})
}
