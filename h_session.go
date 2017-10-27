package webutil

import (
	"crypto/rand"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/noypi/router"
)

type SessionStore struct {
	sstore sessions.Store
}

type _sessionName int

const SessionName _sessionName = 0

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

func CurrentSession(ctx *router.Context) (o *sessions.Session, exists bool) {
	o1, exists := ctx.Get(SessionName)
	return o1.(*sessions.Session), exists
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

func GetSession(ctx *router.Context) *sessions.Session {
	if o, exists := ctx.Get(SessionName); exists {
		return o.(*sessions.Session)
	}

	return nil
}

func (this *SessionStore) addsesion(name string, r *http.Request) (err error) {
	session, err := this.sstore.Get(r, name)
	if nil != err {
		return
	}
	ctx := router.ContextR(r)
	ctx.Set(SessionName, session)
	return
}

func (this *SessionStore) AddSessionHandler(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := this.addsesion(name, r)
		ctx := router.ContextR(r)
		if nil != err {
			w.WriteHeader(http.StatusInternalServerError)
			LogErr(ctx, "AddSessionHandler() err=", err)
		}
	}
}
