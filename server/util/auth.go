package util

import (
	"encoding/gob"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

var sessionName = "session"
var key = []byte("super-secret-key")
var sessionStore = sessions.NewCookieStore(key)

// session values
const (
	sessionUser = "user"
	sessionRole = "role"
	sessionHero = "hero"
)

type RoleLevel int

const (
	RolePlayer RoleLevel = iota
	RoleManager
	RoleAdmin
)

func init() {
	// sessionStore.Options.SameSite = http.SameSiteStrictMode
	sessionStore.Options.SameSite = http.SameSiteLaxMode
	gob.Register(RolePlayer)
}

func AddSession(w http.ResponseWriter, r *http.Request, user string, role RoleLevel) error {
	session, err := sessionStore.Get(r, sessionName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}
	session.Values[sessionUser] = user
	session.Values[sessionRole] = role
	logrus.Infoln("session:", session)
	// 设置session的过期时间
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600, // 比如设置为1小时
		HttpOnly: true, // 增强安全性，防止通过JS访问cookie
	}
	err = session.Save(r, w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	return nil
}

func RemoveSession(w http.ResponseWriter, r *http.Request) error {
	session, err := sessionStore.Get(r, sessionName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}
	session.Options.MaxAge = -1
	session.Save(r, w)
	return nil
}

func GetUser(w http.ResponseWriter, r *http.Request) (user string, role RoleLevel, err error) {
	session, err := sessionStore.Get(r, sessionName)
	logrus.Infoln("session:")
	logrus.Infoln("session:", session)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return "", 0, err
	}
	if session.IsNew {
		w.WriteHeader(http.StatusUnauthorized)
		return "", 0, Errorf("登录信息已失效")
	}
	if _, ok := session.Values[sessionUser]; !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return "", 0, Errorf("登录信息已失效")
	}
	return session.Values[sessionUser].(string),
		session.Values[sessionRole].(RoleLevel),
		nil
}
