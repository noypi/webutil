package webutil

import (
	"context"
	"fmt"
	"net/http"

	"github.com/noypi/router"
)

type _errorKey int

const (
	ErrorKey _errorKey = iota
	ErrorInfo
)

func AddError(ctx context.Context, err error) {
	c := ToStore(ctx)
	c.Set(ErrorKey, err)
}

func GetErrorInfo(ctx context.Context) interface{} {
	c := ToStore(ctx)
	o, _ := c.Get(ErrorInfo)
	return o
}

func ifHasErrorRedirect(w http.ResponseWriter, r *http.Request, theURL string, moreinfo interface{}) {
	ctx := router.ContextW(w)
	c := ToStore(ctx)
	err, hasError := c.Get(ErrorKey)
	if hasError {
		if nil != moreinfo {
			c.Set(ErrorInfo, moreinfo)
		}
		ctx.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s?error=%v", theURL, err))
		ctx.Abort()
		ERR := router.GetErrLog(ctx)
		ERR.Ln("err=", err)
		ERR.PrintStackTrace(20)
		return
	}
}

func IfHasErrorRedirect(theURL string, moreinfo interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ifHasErrorRedirect(w, r, theURL, moreinfo)
	}
}

func OnErrorRedirect(h http.HandlerFunc, theURL string, moreinfo interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
		ifHasErrorRedirect(w, r, theURL, moreinfo)
	}
}

func OnErrorGotoReferrer(h http.HandlerFunc, moreinfo interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
		theURL := r.URL.Path
		if 0 < len(r.URL.RawQuery) {
			theURL += "?" + r.URL.RawQuery
		}
		ifHasErrorRedirect(w, r, theURL, moreinfo)
	}
}
