package webutil

import (
	"net/http"

	cookiejar "github.com/noypi/persistent-cookiejar"
	"github.com/noypi/router"
)

type _persistentCookieKey int

const (
	PersistentCookieKey _persistentCookieKey = iota
)

func GetCookie(ctx *router.Context) *cookiejar.Jar {
	if o, exists := ctx.Get(PersistentCookieKey); exists {
		return o.(*cookiejar.Jar)
	}

	return nil
}

func CreatePersistentCookie(fnGetCookie func(c *router.Context) (*cookiejar.Jar, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := router.ContextR(r)

		jar, err := fnGetCookie(c)
		if nil != err {
			ERR := GetErrLog(c)
			ERR.PrintStackTrace(5)
			ERR.Ln("CreatePersistentCookie: failed to get cookie, err=%s", err.Error())
			AddError(c, err)
			return
		}
		if nil == jar {
			jar, _ = cookiejar.New(nil)
		}

		c.Set(PersistentCookieKey, jar)
	}

}

func SavePersistentCookie(fnSaveCookie func(c *router.Context, jar http.CookieJar) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := router.ContextR(r)
		jar := GetCookie(c)
		fnSaveCookie(c, jar)
	}
}
