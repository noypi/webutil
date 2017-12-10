package webutil

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/gob"
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/noypi/router"
)

type SessionStore struct {
	sstore sessions.Store
}

type _sessionName int

const SessionName _sessionName = 0

func validateOptions(opts *sessions.Options, keys [][]byte) (*sessions.Options, [][]byte) {
	if 0 == len(keys) || 1 == (len(keys)&0x01) {
		// when keys is odd in length, should be in pairs
		keys = GenRandSecrets(3)
	}
	if nil == opts {
		opts = &sessions.Options{
			HttpOnly: true,
			MaxAge:   3600 * 24 * 365 * 100, // (1*365 a year) * 100
		}
	}

	return opts, keys
}

func MarshalKeys(keys [][]byte) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := gob.NewEncoder(buf).Encode(keys)
	return buf.Bytes(), err
}

func UnmarshalKeys(bb []byte) (keys [][]byte, err error) {
	buf := bytes.NewBuffer(bb)
	err = gob.NewEncoder(buf).Encode(&keys)
	return
}

func NewCookieSession(opts *sessions.Options, keys ...[]byte) *SessionStore {
	o := new(SessionStore)
	opts, keys = validateOptions(opts, keys)
	cs := sessions.NewCookieStore(keys...)
	cs.Options = opts
	cs.MaxAge(cs.Options.MaxAge)
	o.sstore = cs
	return o
}

func NewFilesystemSession(fpath string, opts *sessions.Options, keys ...[]byte) *SessionStore {
	o := new(SessionStore)
	opts, keys = validateOptions(opts, keys)
	cs := sessions.NewFilesystemStore(fpath, keys...)
	cs.Options = opts
	cs.MaxAge(cs.Options.MaxAge)
	o.sstore = cs
	return o
}

func CurrentSession(ctx context.Context) (o *sessions.Session, exists bool) {
	c := ToStore(ctx)
	o1, exists := c.Get(SessionName)
	return o1.(*sessions.Session), exists
}

func GenRandSecrets(count int) [][]byte {
	var keys [][]byte
	for i := 0; i < count; i++ {
		bbHashKey := make([]byte, 32)
		rand.Read(bbHashKey)
		keys = append(keys, bbHashKey)

		bbBlockKey := make([]byte, 32)
		rand.Read(bbBlockKey)
		keys = append(keys, bbBlockKey)
	}
	return keys
}

func GetSession(ctx context.Context) *sessions.Session {
	c := ToStore(ctx)
	if o, exists := c.Get(SessionName); exists {
		return o.(*sessions.Session)
	}

	return nil
}

func (this *SessionStore) addsesion(name string, w http.ResponseWriter, r *http.Request) (err error) {
	session, err := this.sstore.Get(r, name)
	if nil != err {
		if nil != session && strings.Contains(err.Error(), "securecookie: the value is not valid") {
			err = nil
		} else {
			// other errors, do abort
			return
		}
	}
	ctx := router.ContextW(w)
	ctx.Set(SessionName, session)
	return
}

func (this *SessionStore) AddSessionHandler(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := router.ContextW(w)
		err := this.addsesion(name, w, r)
		if nil != err {
			ERR := router.GetErrLog(ctx)
			ERR("AddSessionHandler: name=%s, err=%v", name, err)
			AddError(ctx, err)
			return
		}
	}
}

func (this *SessionStore) MaxAge(n int) {
	switch o := this.sstore.(type) {
	case *sessions.CookieStore:
		o.Options.MaxAge = n
		o.MaxAge(n)
	case *sessions.FilesystemStore:
		o.Options.MaxAge = n
		o.MaxAge(n)
	}
}

func (this *SessionStore) Domain(s string) {
	switch o := this.sstore.(type) {
	case *sessions.CookieStore:
		o.Options.Domain = s
	case *sessions.FilesystemStore:
		o.Options.Domain = s
	}
}

func (this *SessionStore) HttpOnly(b bool) {
	switch o := this.sstore.(type) {
	case *sessions.CookieStore:
		o.Options.HttpOnly = b
	case *sessions.FilesystemStore:
		o.Options.HttpOnly = b
	}
}

func (this *SessionStore) Path(s string) {
	switch o := this.sstore.(type) {
	case *sessions.CookieStore:
		o.Options.Path = s
	case *sessions.FilesystemStore:
		o.Options.Path = s
	}
}

func (this *SessionStore) Secure(b bool) {
	switch o := this.sstore.(type) {
	case *sessions.CookieStore:
		o.Options.Secure = b
	case *sessions.FilesystemStore:
		o.Options.Secure = b
	}
}
