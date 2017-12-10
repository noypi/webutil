package webutil

import (
	"html/template"
	"net/http"
)

func UseRenderer(glob string, filenames ...string) http.HandlerFunc {
	var err error
	tpl := template.New("")
	if 0 < len(glob) {
		tpl, err = tpl.ParseGlob(glob)
	} else if 0 < len(filenames) {
		tpl, err = tpl.ParseFiles(filenames...)
	}
	return func(w http.ResponseWriter, r *http.Request) {

		if nil != err {

		}
	}
}
