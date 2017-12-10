package webutil

import (
	"bytes"
	"context"

	"github.com/oxtoacart/bpool"
)

type _bufpooltype int

const (
	BufPoolkey _bufpooltype = iota
)

var BufPoolSize = 16

type BufPool interface {
	Get() *bytes.Buffer
	Put(buf *bytes.Buffer)
}

func MergeMapStrIf(dst, src map[string]interface{}) {
	for k, v := range src {
		dst[k] = v
	}
}

func GetBufPool(ctx context.Context) BufPool {
	c := ToStore(ctx)
	o, exists := c.Get(BufPoolkey)
	if exists {
		return o.(*bpool.BufferPool)
	}
	bufpool := bpool.NewBufferPool(BufPoolSize)
	c.Set(BufPoolkey, bufpool)
	return bufpool
}
