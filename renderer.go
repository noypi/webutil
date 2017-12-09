package webutil

import (
	"context"
	"html/template"
	"io"

	"github.com/noypi/router"
)

type _Renderer struct {
	tpl      *template.Template
	bufpool  BufPool
	c        *router.Context
	tplcache map[string]*template.Template
}

func NewRenderer(c context.Context, t *template.Template) *_Renderer {
	o := new(_Renderer)
	o.tpl = t
	o.c = c.(*router.Context)
	o.bufpool = GetBufPool(o.c)
	o.tplcache = map[string]*template.Template{}
	return o
}

func (r *_Renderer) Render(code int, pages ...interface{}) error {
	if 0 == len(pages) {
		return ErrNoRootTPL
	}

	config0 := GetPageDataKVConfig(r.c, pages[0])

	datamap := map[string]interface{}{}
	MergePagesData(datamap, pages...)

	funcsmap := map[string]interface{}{}
	MergePagesFuncs(funcsmap, pages...)

	tpl := r.CloneTemplate(pages...).Funcs(funcsmap)
	buf := r.bufpool.Get()
	defer r.bufpool.Put(buf)

	err := tpl.ExecuteTemplate(buf, config0["name"], datamap)
	if nil != err {
		return err
	}

	r.c.Writer.WriteHeader(code)
	_, err = io.Copy(r.c.Writer, buf)
	return err
}

func (r *_Renderer) CloneTemplate(pages ...interface{}) *template.Template {
	tpl := template.New("")
	for _, v := range pages {
		config := GetPageDataKVConfig(r.c, v)
		if name, has := config["name"]; has {
			namedtpl := r.tpl.Lookup(name)
			if nil == namedtpl {
				panic("cannot find named template=" + name)
			}
			tpl = template.Must(tpl.AddParseTree(name, namedtpl.Tree.Copy()))
		}
	}
	return tpl
}
