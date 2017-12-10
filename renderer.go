package webutil

import (
	"context"
	"html/template"
	"io"

	"github.com/noypi/router"
)

type _Renderer struct {
	tpl     *template.Template
	bufpool BufPool
	c       *router.Context
}

func NewRenderer(c context.Context, t *template.Template) *_Renderer {
	o := new(_Renderer)
	o.tpl = t
	o.c = c.(*router.Context)
	o.bufpool = GetBufPool(o.c)
	return o
}

func (r *_Renderer) Render(code int, pages ...interface{}) error {
	if 0 == len(pages) {
		return ErrNoRootTPL
	}

	config0 := GetPageDataKVConfig(r.c, pages[0])
	c := ToStore(r.c)

	datamap := map[string]interface{}{}
	MergePagesData(c, datamap, pages...)

	funcsmap := map[string]interface{}{}
	MergePagesFuncs(c, funcsmap, pages...)

	tpl, err := r.CloneTemplate(pages...)
	if nil != err {
		return err
	}
	tpl = tpl.Funcs(funcsmap)
	buf := r.bufpool.Get()
	defer r.bufpool.Put(buf)

	if err = tpl.ExecuteTemplate(buf, config0["name"], datamap); nil != err {
		return err
	}

	r.c.Writer.WriteHeader(code)
	_, err = io.Copy(r.c.Writer, buf)
	return err
}

func (r *_Renderer) CloneTemplate(pages ...interface{}) (tpl *template.Template, err error) {
	tpl = template.New("")
	for _, v := range pages {
		config := GetPageDataKVConfig(r.c, v)
		name, has := config["name"]
		if !has {
			continue
		}

		namedtpl := r.tpl.Lookup(name)
		if nil == namedtpl {
			if o, has := v.(HasHTML); has {
				if err = addHTMLContent(tpl, name, o); nil != err {
					return
				}
			} else {
				panic("cannot find named template=" + name)
			}

		} else {
			template.Must(tpl.AddParseTree(name, namedtpl.Tree))
		}

	}
	return tpl, nil
}

func addHTMLContent(tpl *template.Template, name string, o HasHTML) error {
	if content := o.GetHTMLContent(); 0 < len(content) {
		if temp, err := template.New(name).Parse(content); nil != err {
			return err
		} else {
			tpl.AddParseTree(name, temp.Tree)
		}
	}
	return nil
}
