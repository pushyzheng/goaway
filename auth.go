package main

import (
	"fmt"
	logger "github.com/sirupsen/logrus"
	"github.com/snluu/uuid"
	"html/template"
	"net/http"
	"sync"
	"time"
)

type User struct {
	SessionId string    `json:"sessionId"`
	Username  string    `json:"username"`
	Expires   time.Time `json:"expires"`
}

var sessions = sync.Map{}

// Login Return login html page
func Login(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "text/html; charset=utf-8")
	t, err := template.ParseFiles("static/login.html")
	if err != nil {
		logger.Errorln("parse file error:", err.Error())
		_, _ = fmt.Fprintf(resp, "Unable to load template")
		return
	}
	_ = t.Execute(resp, nil)
}

// Submit parse form, then check username and password
func Submit(resp http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		logger.Errorln("parse form error:", err.Error())
		_, _ = fmt.Fprintf(resp, "Fail to parse form: %s", err.Error())
		return
	}
	username := req.FormValue("username")
	if _, ok := Conf.Accounts[username]; !ok {
		ReturnError(resp, http.StatusUnauthorized, "The account is unavailable")
		return
	}
	account := Conf.Accounts[username]
	password := req.FormValue("password")
	if password != account.Password {
		ReturnError(resp, http.StatusUnauthorized, "Invalid username or password")
		return
	}
	// Login succeed, set cookie to client and save session
	id := uuid.Rand().Hex()
	expires := time.Now().Add(Conf.Server.CookieExpiredHours * time.Hour)
	http.SetCookie(resp, &http.Cookie{
		Name:    IdentityKeyName,
		Value:   id,
		Path:    "/",
		Expires: expires,
		Domain:  Conf.Server.Domain,
		MaxAge:  90000,
	})
	sessions.Store(id, User{SessionId: id, Username: username, Expires: expires})
	http.Redirect(resp, req, "/", http.StatusSeeOther)
}

func HasLogin(r *http.Request) bool {
	_, ok := GetUser(r)
	return ok
}

func GetUser(r *http.Request) (User, bool) {
	cookie, _ := r.Cookie(IdentityKeyName)
	if cookie == nil {
		return User{}, false
	}
	user, ok := sessions.Load(cookie.Value)
	if !ok {
		return User{}, false
	}
	return user.(User), true
}
