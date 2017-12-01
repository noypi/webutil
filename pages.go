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

func GetPageDataKVConfig(ctx context.Context, page interface{}) (m map[string]string) {
	c := ctx.(*router.Context)
	m = map[string]string{}

	// get tplname
	t := reflect.TypeOf(page)
	cachek := fmt.Sprintf("%s,%s", t.PkgPath(), t.String())
	cache := GetPageTypeCache(c)
	if o, has := cache[cachek]; has {
		return o
	}

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

	cache[cachek] = m

	return

}

func MergePagesData(datamap map[string]interface{}, pages ...interface{}) {
	for i := len(pages) - 1; 0 <= i; i-- {
		//MergeMapStrIf(datamap, pages[i].MapData())
		MergePageData(datamap, pages[i])
	}
	return
}

func MergePageData(m map[string]interface{}, page interface{}) {
	v := reflect.ValueOf(page)
	for reflect.Ptr == v.Kind() {
		v = reflect.Indirect(v)
	}
	m[v.Type().Name()] = page
}

// example of page[0] struct
// type SomePage struct {
//		TPLConfig string `webtpl:"name=mytemplate.tpl"`
//	}
func RenderPage(ctx context.Context, code int, pages ...interface{}) {
	if 0 == len(pages) {
		return
	}

	c := ctx.(*router.Context)
	config0 := GetPageDataKVConfig(c, pages[0])

	datamap := map[string]interface{}{}
	MergePagesData(datamap, pages...)

	c.HTML(code, config0["name"], datamap)
}
