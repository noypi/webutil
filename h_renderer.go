package webutil

import (
	"context"
	"html/template"
	"net/http"

	"github.com/noypi/router"
)

type _rendererKey int

const (
	RendererKey _rendererKey = iota
)

func GetRenderer(ctx context.Context) *_Renderer {
	c := ToStore(ctx)
	if o, has := c.Get(RendererKey); has {
		return o.(*_Renderer)
	}
	return nil
}

func UseRenderer(glob string, filenames ...string) http.HandlerFunc {
	var err error
	tpl := template.New("")
	if 0 < len(glob) {
		tpl, err = tpl.ParseGlob(glob)
	} else if 0 < len(filenames) {
		tpl, err = tpl.ParseFiles(filenames...)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		c := router.ContextW(w)
		if nil != err {
			AddError(c, err)
			c.Abort()
			return
		}
		renderer := NewRenderer(c, tpl)
		c.Set(RendererKey, renderer)
	}
}
