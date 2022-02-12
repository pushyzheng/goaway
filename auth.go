package main

import (
	"fmt"
	"github.com/snluu/uuid"
	"html/template"
	"log"
	"net/http"
	"time"
)

var Sessions = make(map[string]string)

func Login(w http.ResponseWriter) {
	// Return login html page
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t, err := template.ParseFiles("login.html")
	if err != nil {
		log.Printf("parse file error: %s\n", err.Error())
		_, _ = fmt.Fprintf(w, "Unable to load template")
		return
	}
	t.Execute(w, nil)
}

func Submit(w http.ResponseWriter, r *http.Request) {
	// parse form, then check username and password
	err := r.ParseForm()
	if err != nil {
		log.Printf("parse form error: %s\n", err.Error())
		_, _ = fmt.Fprintf(w, "Fail to parse form: %s", err.Error())
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	if username != conf.Username || password != conf.Password {
		returnError(w, http.StatusUnauthorized, "INVALID USERNAME OR PASSWORD")
		return
	}
	// Login succeed, set cookie to client and save session
	id := uuid.Rand().Hex()
	http.SetCookie(w, &http.Cookie{
		Name:    IdentityKeyName,
		Value:   id,
		Path:    "/",
		Expires: time.Now().Add(conf.CookieExpiredHours * time.Hour),
		Domain:  conf.Domain,
		MaxAge:  90000,
	})
	Sessions[id] = id
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
