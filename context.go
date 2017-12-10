package webutil

import (
	"context"
)

type Store interface {
	Set(k, v interface{})
	Get(k interface{}) (interface{}, bool)
}

func ToStore(c context.Context) Store {
	o, ok := c.(Store)
	if !ok {
		panic("unknown context type")
	}
	return o
}
