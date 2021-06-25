package main

import (
	"fmt"
	"github.com/snluu/uuid"
	"html/template"
	"net/http"
	"time"
)

var Sessions = make(map[string]string)

func Login(w http.ResponseWriter) {
	// return Login.html page
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t, err := template.ParseFiles("login.html")
	if err != nil {
		fmt.Fprintf(w, "Unable to load template")
	}
	t.Execute(w, nil)
}

func Submit(w http.ResponseWriter, r *http.Request) {
	// parse form, then check username and password
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	if username != "admin" || password != "123456" {
		http.Redirect(w, r, "/login/", http.StatusSeeOther)
		return
	}
	// Login succeed, set cookie to client and save session
	id := string(uuid.Rand().Hex())
	http.SetCookie(w, &http.Cookie{
		Name:    "SESSIONID",
		Value:   id,
		Path:    "/",
		Expires: time.Now().Add(10 * time.Minute),
		Domain:  "pushyzheng.com",
		MaxAge:  90000,
	})
	Sessions[id] = id
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
