package webutil

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/noypi/router"
)

const TPLFieldConfig = "TPLConfig"
const TPLName = "webtpl"

type _pageTypeCacheKey int

const (
	PageTypeCache _pageTypeCacheKey = iota
)

type HasHTML interface {
	GetHTMLContent() string
}

func GetPageTypeCache(ctx context.Context) (m map[string]map[string]string) {
	c := ctx.(*router.Context)
	o, has := c.Get(PageTypeCache)
	if !has {
		m = map[string]map[string]string{}
		c.Set(PageTypeCache, m)
	} else {
		m = o.(map[string]map[string]string)
	}
	return
}

// example of page[0] struct
// type SomePage struct {
//		TPLConfig string `webtpl:"name=mytemplate.tpl"`
//	}
func GetPageDataKVConfig(ctx context.Context, page interface{}) (m map[string]string) {
	var c *router.Context
	if nil != ctx {
		c = ctx.(*router.Context)
	}

	// get tplname
	t := reflect.TypeOf(page)

	// use cache
	var cachek string
	var cache map[string]map[string]string
	if nil != c {
		cachek = fmt.Sprintf("%s,%s", t.PkgPath(), t.String())
		cache = GetPageTypeCache(c)
		if o, has := cache[cachek]; has {
			return o
		}
	}

	m = map[string]string{}
	f, b := t.FieldByName(TPLFieldConfig)
	if !b {
		return
	}

	tag := f.Tag.Get(TPLName)
	if 0 == len(tag) {
		return
	}

	for _, kv := range strings.Split(tag, ",") {
		pair := strings.Split(kv, "=")
		if 2 == len(pair) {
			m[pair[0]] = pair[1]
		}
	}

	if nil != c {
		cache[cachek] = m
	}

	return

}

func MergePagesData(datamap map[string]interface{}, pages ...interface{}) {
	for i := len(pages) - 1; 0 <= i; i-- {
		MergePageData(datamap, pages[i])
	}
	return
}

func MergePagesFuncs(funcmap map[string]interface{}, pages ...interface{}) {
	for i := len(pages) - 1; 0 <= i; i-- {
		MergePageFuncs(funcmap, pages[i])
	}
	return
}

func MergePageData(data map[string]interface{}, page interface{}) {
	v := reflect.ValueOf(page)
	for reflect.Ptr == v.Kind() {
		v = reflect.Indirect(v)
	}
	t := v.Type()
	name := t.Name()

	data["o"+name] = page
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if f.CanInterface() {
			if f.Kind() == reflect.Struct {
				MergePageData(data, f.Interface())
			} else {
				data[t.Field(i).Name] = f.Interface()
			}

		}
	}
}

func MergePageFuncs(funcs map[string]interface{}, page interface{}) {
	v := reflect.ValueOf(page)
	for reflect.Ptr == v.Kind() {
		v = reflect.Indirect(v)
	}
	t := v.Type()
	name := t.Name()
	for i := 0; i < v.NumMethod(); i++ {
		f := v.Method(i)
		if f.CanInterface() {
			funcs[name+t.Method(i).Name] = f.Interface()
		}
	}
}
