package webutil

import (
	"log"
	"net/http"

	"github.com/noypi/logfn"
	"github.com/noypi/router"
)

type _logFuncType int

const (
	LogErrName _logFuncType = iota
	LogInfoName
	LogWarnName
)

func WithErrLogger(ctx *router.Context, fn logfn.LogFunc) *router.Context {
	ctx.Set(LogErrName, fn)
	return ctx
}

func WithWarnLogger(ctx *router.Context, fn logfn.LogFunc) *router.Context {
	ctx.Set(LogWarnName, fn)
	return ctx
}

func WithInfoLogger(ctx *router.Context, fn logfn.LogFunc) *router.Context {
	ctx.Set(LogInfoName, fn)
	return ctx
}

func GetErrLog(ctx *router.Context) logfn.LogFunc {
	return getLogFunc(ctx, LogErrName)
}

func GetInfoLog(ctx *router.Context) logfn.LogFunc {
	return getLogFunc(ctx, LogInfoName)
}

func GetWarnLog(ctx *router.Context) logfn.LogFunc {
	return getLogFunc(ctx, LogWarnName)
}

func getLogFunc(ctx *router.Context, name _logFuncType) (fn logfn.LogFunc) {
	if nil == ctx {
		return log.Printf
	}

	if o, exists := ctx.Get(name); exists {
		fn = (o).(logfn.LogFunc)
	} else {
		fn = log.Printf
	}
	return
}

func LogErr(ctx *router.Context, fmt string, params ...interface{}) {
	GetErrLog(ctx)(fmt, params...)
}

func LogInfo(ctx *router.Context, fmt string, params ...interface{}) {
	GetInfoLog(ctx)(fmt, params...)
}

func LogWarn(ctx *router.Context, fmt string, params ...interface{}) {
	GetWarnLog(ctx)(fmt, params...)
}

func AddLoggerHandler(fnInfo, fnErr, fnWarn logfn.LogFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := router.ContextR(r)
		if nil != fnInfo {
			WithInfoLogger(ctx, fnInfo)
		} else if nil != fnErr {
			WithErrLogger(ctx, fnErr)
		} else if nil != fnWarn {
			WithWarnLogger(ctx, fnWarn)
		}
	})
}
func ErrLoggerHandler(fn logfn.LogFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := router.ContextR(r)
		WithErrLogger(ctx, fn)
	})
}

func InfoLoggerHandler(fn logfn.LogFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WithInfoLogger(router.ContextR(r), fn)
	})
}
