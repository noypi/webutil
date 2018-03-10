package webutil

import (
	"context"
	"html/template"
	"io"
	"strings"

	"github.com/noypi/router"
)

type _Renderer struct {
	tpl     *template.Template
	bufpool BufPool
	c       context.Context
}

func NewRenderer(c context.Context, t *template.Template) *_Renderer {
	o := new(_Renderer)
	o.tpl = t
	o.c = c
	o.bufpool = GetBufPool(o.c)
	return o
}

func (r *_Renderer) Execute(code int, name string, data interface{}) error {
	namedtpl := r.tpl.Lookup(name)
	if nil == namedtpl {
		return ErrNoRootTPL
	}

	buf := r.bufpool.Get()
	defer r.bufpool.Put(buf)

	writer := router.GetWriter(r.c)
	writer.WriteHeader(code)
	err := template.Must(namedtpl.Clone()).ExecuteTemplate(buf, name, data)
	if nil != err {
		return err
	}

	_, err = io.Copy(writer, buf)
	return err
}

func (r *_Renderer) Render(code int, pages ...interface{}) error {
	if 0 == len(pages) {
		return ErrNoRootTPL
	}

	config0 := GetPageDataKVConfig(ToStore(r.c), pages[0])
	c := ToStore(r.c)

	datamap := map[string]interface{}{}
	MergePagesData(c, datamap, pages...)

	funcsmap := map[string]interface{}{}
	MergePagesFuncs(c, funcsmap, pages...)

	DBG := router.GetDebugLog(r.c)
	DBG.Ln("Render datmap=", datamap)

	tpl, err := r.CloneTemplate(pages...)
	if nil != err {
		return err
	}
	tpl = tpl.Funcs(funcsmap)
	buf := r.bufpool.Get()
	defer r.bufpool.Put(buf)

	DBG.Ln("Render fields 0=", strings.Fields(config0["name"])[0])

	if err = tpl.ExecuteTemplate(buf, strings.Fields(config0["name"])[0], datamap); nil != err {
		return err
	}

	writer := router.GetWriter(r.c)
	writer.WriteHeader(code)
	_, err = io.Copy(writer, buf)
	return err
}

func (r *_Renderer) CloneTemplate(pages ...interface{}) (tpl *template.Template, err error) {
	tpl = template.New("")
	for _, v := range pages {
		config := GetPageDataKVConfig(ToStore(r.c), v)
		sname, has := config["name"]
		if !has {
			continue
		}

		names := strings.Fields(sname)
		for _, name := range names {
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
