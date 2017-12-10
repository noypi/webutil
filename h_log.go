package webutil

import (
	"log"
	"net/http"

	"context"

	"github.com/noypi/logfn"
	"github.com/noypi/router"
)

type _logFuncType int

const (
	LogErrName _logFuncType = iota
	LogInfoName
	LogWarnName
	LogDebugName
)

func WithErrLogger(ctx context.Context, fn logfn.LogFunc) context.Context {
	c := ToStore(ctx)
	c.Set(LogErrName, fn)
	return ctx
}

func WithWarnLogger(ctx context.Context, fn logfn.LogFunc) context.Context {
	c := ToStore(ctx)
	c.Set(LogWarnName, fn)
	return ctx
}

func WithInfoLogger(ctx context.Context, fn logfn.LogFunc) context.Context {
	c := ToStore(ctx)
	c.Set(LogInfoName, fn)
	return ctx
}

func WithDebugLogger(ctx context.Context, fn logfn.LogFunc) context.Context {
	c := ToStore(ctx)
	c.Set(LogDebugName, fn)
	return ctx
}

func GetErrLog(ctx context.Context) logfn.LogFunc {
	return getLogFunc(ctx, LogErrName)
}

func GetInfoLog(ctx context.Context) logfn.LogFunc {
	return getLogFunc(ctx, LogInfoName)
}

func GetWarnLog(ctx context.Context) logfn.LogFunc {
	return getLogFunc(ctx, LogWarnName)
}

func GetDebugLog(ctx context.Context) logfn.LogFunc {
	return getLogFunc(ctx, LogDebugName)
}

func getLogFunc(ctx context.Context, name _logFuncType) (fn logfn.LogFunc) {
	if nil == ctx {
		return log.Printf
	}

	c := ToStore(ctx)

	if o, exists := c.Get(name); exists {
		fn = (o).(logfn.LogFunc)
	} else {
		fn = log.Printf
	}
	return
}

func LogErr(ctx context.Context, fmt string, params ...interface{}) {
	GetErrLog(ctx)(fmt, params...)
}

func LogInfo(ctx context.Context, fmt string, params ...interface{}) {
	GetInfoLog(ctx)(fmt, params...)
}

func LogWarn(ctx context.Context, fmt string, params ...interface{}) {
	GetWarnLog(ctx)(fmt, params...)
}

func LogDebug(ctx context.Context, fmt string, params ...interface{}) {
	GetDebugLog(ctx)(fmt, params...)
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
