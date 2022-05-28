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

const (
	defaultExpiredHoursOfCookies = 24
)

type User struct {
	SessionId string    `json:"sessionId"`
	Username  string    `json:"username"`
	Expires   time.Time `json:"expires"`
}

var sessions = sync.Map{}

// Login Return login html page
func Login(resp http.ResponseWriter, req *http.Request) {
	if HasLogin(req) {
		http.Redirect(resp, req, "/", http.StatusSeeOther)
		return
	}
	resp.Header().Set("Content-Type", "text/html; charset=utf-8")
	t, err := template.ParseFiles(LoginPagePath)
	if err != nil {
		logger.Errorln("parse file error:", err.Error())
		_, _ = fmt.Fprintf(resp, "Unable to load template")
		return
	}
	err = t.Execute(resp, nil)
	if err != nil {
		logger.Errorf("template execute error: %s %s", LoginPagePath, err)
	}
}

// Submit parse form, then check username and password
func Submit(resp http.ResponseWriter, req *http.Request) {
	if HasLogin(req) {
		http.Redirect(resp, req, "/", http.StatusSeeOther)
		return
	}
	err := req.ParseForm()
	if err != nil {
		logger.Errorln("parse form error:", err.Error())
		ReturnError(resp, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	username := req.FormValue("username")
	if _, ok := Conf.Accounts[username]; !ok {
		ReturnError(resp, http.StatusUnauthorized, "The account is unavailable")
		return
	}
	account := Conf.Accounts[username]
	pwd := req.FormValue("password")
	if pwd != account.Password {
		ReturnError(resp, http.StatusUnauthorized, "Invalid username or password")
		return
	}
	// Login succeed, set cookie to client and save session
	id := uuid.Rand().Hex()
	expired := Conf.Server.CookieExpiredHours
	if expired == 0 {
		expired = defaultExpiredHoursOfCookies
	}
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

func Logout(resp http.ResponseWriter, req *http.Request) {
	if user, ok := GetUser(req); ok {
		sessions.Delete(user.SessionId)
	}
	DirectLogin(resp, req)
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
