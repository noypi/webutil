package webutil

import (
	"crypto/rand"
	"net/http"

	"context"

	"github.com/gorilla/sessions"
)

var g_sessionsCookieStore *sessions.CookieStore
var g_sessionsFilesystemStore *sessions.FilesystemStore
var g_sessionsSecrets [][]byte

const (
	SessionName = "$session"
)

func SetSessionsSecret(keys ...[]byte) {
	for _, bb := range keys {
		g_sessionsSecrets = append(g_sessionsSecrets, bb)
	}
}

func genRandSecrets() {
	for i := 0; i < 3; i++ {
		bb := make([]byte, 10)
		rand.Read(bb)
		g_sessionsSecrets = append(g_sessionsSecrets, bb)
	}
}

func addsesion(sstore sessions.Store, name string, r *http.Request) (ctx context.Context, err error) {
	session, err := sstore.Get(r, name)
	if nil != err {
		return
	}
	ctx = r.Context()
	return context.WithValue(ctx, SessionName, session), nil
}

func AddCookieSession(name string, nexth http.Handler) http.Handler {
	if 0 == len(g_sessionsSecrets) {
		genRandSecrets()
	}
	if nil == g_sessionsCookieStore {
		g_sessionsCookieStore = sessions.NewCookieStore(g_sessionsSecrets...)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, err := addsesion(g_sessionsCookieStore, name, r)
		if nil != err {
			w.WriteHeader(http.StatusInternalServerError)
			LogErr(ctx, "AddCookieSession() err=", err)
		} else {
			nexth.ServeHTTP(w, r.WithContext(ctx))
		}

	})
}

func AddFilesystemSession(path, name string, nexth http.Handler) http.Handler {
	if 0 == len(g_sessionsSecrets) {
		genRandSecrets()
	}
	if nil == g_sessionsFilesystemStore {
		g_sessionsFilesystemStore = sessions.NewFilesystemStore(path, g_sessionsSecrets...)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, err := addsesion(g_sessionsFilesystemStore, name, r)
		if nil != err {
			w.WriteHeader(http.StatusInternalServerError)
			LogErr(ctx, "AddFilesystemSession() err=", err)
		} else {
			nexth.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}
