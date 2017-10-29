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

func NewCookieSession(domain string, keys ...[]byte) *SessionStore {
	if 0 == len(keys) || 1 == (len(keys)&0x01) {
		keys = genRandSecrets()
	}

	o := new(SessionStore)
	cs := sessions.NewCookieStore(keys...)
	cs.Options = &sessions.Options{
		//Path:     "/",
		MaxAge: 3600 * 24 * 365 * 100, // (1*365 a year) * 100
		//Secure:   true,
		HttpOnly: true,
	}
	if 0 < len(domain) {
		cs.Options.Domain = domain
	}
	cs.MaxAge(cs.Options.MaxAge)
	o.sstore = cs

	return o
}

func NewFilesystemSession(domain string, path string, keys ...[]byte) *SessionStore {
	if 0 == len(keys) || 1 == (len(keys)&0x01) {
		keys = genRandSecrets()
	}

	o := new(SessionStore)
	cs := sessions.NewFilesystemStore(path, keys...)
	cs.Options = &sessions.Options{
		//Path:     "/",
		MaxAge: 3600 * 24 * 365 * 100, // (1*365 a year) * 100
		//Secure:   true,
		HttpOnly: true,
	}
	if 0 < len(domain) {
		cs.Options.Domain = domain
	}
	cs.MaxAge(cs.Options.MaxAge)
	o.sstore = sessions.NewFilesystemStore(path, keys...)

	return o
}

func CurrentSession(ctx *router.Context) (o *sessions.Session, exists bool) {
	o1, exists := ctx.Get(SessionName)
	return o1.(*sessions.Session), exists
}

func genRandSecrets() [][]byte {
	var keys [][]byte
	for i := 0; i < 1; i++ {
		bbHashKey := make([]byte, 32)
		rand.Read(bbHashKey)
		keys = append(keys, bbHashKey)

		bbBlockKey := make([]byte, 32)
		rand.Read(bbBlockKey)
		keys = append(keys, bbBlockKey)
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
		ctx := router.ContextR(r)
		err := this.addsesion(name, r)
		if nil != err {
			LogErr(ctx, "AddSessionHandler() err=%v, name=%s", err, name)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}
}
