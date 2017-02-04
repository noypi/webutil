package webutil

import (
	"crypto/rand"
	"net/http"

	"context"

	"github.com/gorilla/sessions"
	"github.com/noypi/util"
)

type SessionStore struct {
	sstore sessions.Store
}

const (
	SessionName = "$session"
)

func NewCookieSession(keys ...[]byte) *SessionStore {
	if 0 == len(keys) {
		keys = genRandSecrets()
	}
	o := new(SessionStore)
	o.sstore = sessions.NewCookieStore(keys...)

	return o
}

func NewFilesystemSession(path string, keys ...[]byte) *SessionStore {
	if 0 == len(keys) {
		keys = genRandSecrets()
	}
	o := new(SessionStore)
	o.sstore = sessions.NewFilesystemStore(path, keys...)

	return o
}

func CurrentSession(ctx context.Context) *sessions.Session {
	return ctx.Value(SessionName).(*sessions.Session)
}

func genRandSecrets() [][]byte {
	var keys [][]byte
	for i := 0; i < 3; i++ {
		bb := make([]byte, 10)
		rand.Read(bb)
		keys = append(keys, bb)
	}
	return keys
}

func (this *SessionStore) addsesion(name string, r *http.Request) (ctx context.Context, err error) {
	session, err := this.sstore.Get(r, name)
	if nil != err {
		return
	}
	ctx = r.Context()
	return context.WithValue(ctx, SessionName, session), nil
}

func (this *SessionStore) AddSessionHandler(name string, nexth http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, err := this.addsesion(name, r)
		if nil != err {
			w.WriteHeader(http.StatusInternalServerError)
			util.LogErr(ctx, "AddSessionHandler() err=", err)
		} else {
			nexth.ServeHTTP(w, r.WithContext(ctx))
		}

	})
}
