package webutil

import (
	"net/http"

	"github.com/go-zoo/bone"
)

type MethodType int

const (
	MethodGet MethodType = iota
	MethodPost
	MethodDelete
	MethodPut
	MethodPatch
	MethodHead
	MethodOptions
)

type Mux interface {
	Handle(string, http.Handler)
	HandleFunc(string, http.HandlerFunc)
	ServeHTTP(http.ResponseWriter, *http.Request)
}

func NewPreferredMux() Mux {
	return bone.New()
}

func SetNotFound(mux Mux, h http.Handler) {
	mux.(*bone.Mux).NotFound(h)
}

func SetRoutesFunc(fn func(string, http.HandlerFunc), m map[string]http.HandlerFunc) {
	for k, v := range m {
		fn(k, v)
	}
}

func SetRoutes(fn func(string, http.Handler), m map[string]http.Handler) {
	for k, v := range m {
		fn(k, v)
	}
}

func HandlerMethod(mux Mux, method MethodType) func(string, http.Handler) {
	fn := muxBoneHandlerMethod(mux, method)
	return func(s string, h http.Handler) {
		fn(s, h)
	}
}

func HandlerFuncMethod(mux Mux, method MethodType) func(string, http.HandlerFunc) {
	fn := muxBoneHandlerFuncMethod(mux, method)
	return func(s string, h http.HandlerFunc) {
		fn(s, h)
	}
}

func muxBoneHandlerMethod(mux Mux, method MethodType) (fn func(string, http.Handler) (fn *bone.Route)) {
	switch method {
	case MethodGet:
		fn = mux.(*bone.Mux).Get
	case MethodPost:
		fn = mux.(*bone.Mux).Post
	case MethodDelete:
		fn = mux.(*bone.Mux).Delete
	case MethodPut:
		fn = mux.(*bone.Mux).Put
	case MethodPatch:
		fn = mux.(*bone.Mux).Patch
	case MethodHead:
		fn = mux.(*bone.Mux).Head
	case MethodOptions:
		fn = mux.(*bone.Mux).Options
	}
	return fn
}

func muxBoneHandlerFuncMethod(mux Mux, method MethodType) (fn func(string, http.HandlerFunc) *bone.Route) {
	switch method {
	case MethodGet:
		fn = mux.(*bone.Mux).GetFunc
	case MethodPost:
		fn = mux.(*bone.Mux).PostFunc
	case MethodDelete:
		fn = mux.(*bone.Mux).DeleteFunc
	case MethodPut:
		fn = mux.(*bone.Mux).PutFunc
	case MethodPatch:
		fn = mux.(*bone.Mux).PatchFunc
	case MethodHead:
		fn = mux.(*bone.Mux).HeadFunc
	case MethodOptions:
		fn = mux.(*bone.Mux).OptionsFunc
	}
	return fn
}
