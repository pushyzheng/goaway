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

var sessions = sync.Map{}

// Login Return login html page
func Login(resp http.ResponseWriter) {
	resp.Header().Set("Content-Type", "text/html; charset=utf-8")
	t, err := template.ParseFiles("static/login.html")
	if err != nil {
		logger.Errorln("parse file error:", err.Error())
		_, _ = fmt.Fprintf(resp, "Unable to load template")
		return
	}
	t.Execute(resp, nil)
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
	password := req.FormValue("password")
	if username != conf.Username || password != conf.Password {
		returnError(resp, http.StatusUnauthorized, "INVALID USERNAME OR PASSWORD")
		return
	}
	// Login succeed, set cookie to client and save session
	id := uuid.Rand().Hex()
	http.SetCookie(resp, &http.Cookie{
		Name:    IdentityKeyName,
		Value:   id,
		Path:    "/",
		Expires: time.Now().Add(conf.CookieExpiredHours * time.Hour),
		Domain:  conf.Domain,
		MaxAge:  90000,
	})
	sessions.Store(id, id)
	http.Redirect(resp, req, "/", http.StatusSeeOther)
}

func HasLogin(r *http.Request) bool {
	cookie, _ := r.Cookie(IdentityKeyName)
	if cookie == nil {
		return false
	}
	_, ok := sessions.Load(cookie.Value)
	return ok
}
