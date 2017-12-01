package webutil

import (
	"context"
	"fmt"
	"html/template"
	"net/http"

	"github.com/noypi/router"
)

type _templateType int

const TPLRootKey _templateType = 0

var ErrNoRootTPL = fmt.Errorf("no root template found.")

func GetRootTPL(ctx context.Context) *template.Template {
	c := ctx.(*router.Context)
	if t, b := c.Get(TPLRootKey); b {
		return t.(*template.Template)
	}
	return nil
}

func SetRootTemplate(ctx context.Context, tpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := ctx.(*router.Context)
		c.Set(TPLRootKey, tpl)
	}
}

func AddGlobToRootTPL(pattern string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := router.ContextR(r)
		roottpl := GetRootTPL(c)
		if nil == roottpl {
			AddError(c, ErrNoRootTPL)
			return
		}

		if _, err := roottpl.ParseGlob(pattern); nil != err {
			AddError(c, err)
			return
		}
	}
}

func AddFilesToRootTPL(filenames ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := router.ContextR(r)
		roottpl := GetRootTPL(c)
		if nil == roottpl {
			AddError(c, ErrNoRootTPL)
			return
		}

		if _, err := roottpl.ParseFiles(filenames...); nil != err {
			AddError(c, err)
			return
		}
	}
}
