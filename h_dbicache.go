package webutil

import (
	"net/http"

	"reflect"

	"bitbucket.org/noypi/gae"
	"github.com/noypi/router"
)

func UseDbiCache(namespace string, dbiKey interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := router.ContextW(w)
		oDbi, exists := c.Get(dbiKey)
		if !exists {
			return
		}

		oDbiCache := GetDbiCache(namespace, oDbi.(gae.DbExInt))
		c.Set(dbiKey, oDbiCache)
	}
}

type _cachetype map[string]interface{}

type DbiCache struct {
	gae.DbExInt
	cache _cachetype
}

var g_cache = map[string]_cachetype{}

func getCache(namespace string) _cachetype {
	m, _ := g_cache[namespace]
	if nil == m {
		m = _cachetype{}
		g_cache[namespace] = m
	}
	return m
}

func GetDbiCache(namespace string, oDbi gae.DbExInt) *DbiCache {
	o := new(DbiCache)
	o.DbExInt = oDbi
	o.cache = getCache(namespace)
	return o
}

func (this *DbiCache) Put(k string, bb []byte) error {
	this.cache[k] = bb
	return this.DbExInt.Put(k, bb)
}

func (this *DbiCache) Get(k string) (bb []byte, err error) {
	bb, has := this.cache[k].([]byte)
	if !has {
		bb, err = this.DbExInt.Get(k)
	}
	return
}

func (this *DbiCache) PutObject(k string, o interface{}) error {
	this.cache[k] = o
	return this.DbExInt.PutObject(k, o)
}

func (this *DbiCache) GetObject(k string, o interface{}) (err error) {
	tmp, has := this.cache[k]
	if !has {
		err = this.DbExInt.GetObject(k, o)
	} else {
		v1 := reflect.Indirect(reflect.ValueOf(o))
		v2 := reflect.Indirect(reflect.ValueOf(tmp))
		for i := 0; i < v1.NumField(); i++ {
			if v1.Field(i).CanSet() {
				v1.Field(i).Set(v2.Field(i))
			}
		}
	}
	return
}
