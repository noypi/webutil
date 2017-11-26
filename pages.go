package webutil

import (
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

type IPage interface {
	MapData() map[string]interface{}
}

type Pages []IPage

func GetPageTypeCache(c *router.Context) (m map[string]map[string]string) {
	o, has := c.Get(PageTypeCache)
	if !has {
		m = map[string]map[string]string{}
		c.Set(PageTypeCache, m)
	} else {
		m = o.(map[string]map[string]string)
	}
	return
}

func GetPageDataKVConfig(c *router.Context, page interface{}) (m map[string]string) {
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

func (pages Pages) MergeData(datamap map[string]interface{}) {
	for i := len(pages) - 1; 0 <= i; i-- {
		// merge data
		MergeMapStrIf(datamap, pages[i].MapData())
	}
	return
}

// example of page[0] struct
// type SomePage struct {
//		TPLConfig string `webtpl:"name=mytemplate.tpl"`
//	}
func RenderPage(c *router.Context, code int, pages ...IPage) {
	if 0 == len(pages) {
		return
	}

	config0 := GetPageDataKVConfig(c, pages[0])

	ps := Pages(pages)
	datamap := map[string]interface{}{}
	ps.MergeData(datamap)

	c.HTML(code, config0["name"], datamap)
}
