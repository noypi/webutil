package webutil

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/noypi/logfn"
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

func GetPageTypeCache(c Store) (m map[string]map[string]string) {
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
func GetPageDataKVConfig(c Store, page interface{}) (m map[string]string) {
	// get tplname
	t := reflect.TypeOf(page)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

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

func MergePagesData(c Store, datamap map[string]interface{}, pages ...interface{}) {
	for i := len(pages) - 1; 0 <= i; i-- {
		MergePageData(c, datamap, pages[i])
	}
	return
}

func MergePagesFuncs(c Store, funcmap map[string]interface{}, pages ...interface{}) {
	for i := len(pages) - 1; 0 <= i; i-- {
		MergePageFuncs(c, funcmap, pages[i])
	}
	return
}

func MergePageData(c Store, data map[string]interface{}, page interface{}) {
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
				MergePageData(c, data, f.Interface())
			} else {
				data[t.Field(i).Name] = f.Interface()
			}

		}
	}
}

func MergePageFuncs(c Store, funcs map[string]interface{}, page interface{}) {
	o, _ := c.Get(router.LogDebugName)
	DBG := o.(logfn.LogFunc)
	DBG.Ln("+***********************************MergePageFuncs...")
	defer DBG.Ln("-***********************************MergePageFuncs...")

	v := reflect.ValueOf(page)
	for reflect.Ptr == v.Kind() {
		v = reflect.Indirect(v)
	}
	t := v.Type()
	name := t.Name()

	for i := 0; i < v.NumMethod(); i++ {
		f := v.Method(i)

		if f.CanInterface() {
			DBG.Ln("MergePageFuncs...k=", name+t.Method(i).Name)
			funcs[name+t.Method(i).Name] = f.Interface()
		}
	}
}
