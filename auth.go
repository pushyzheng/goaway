package main

import (
	"fmt"
	"github.com/snluu/uuid"
	"html/template"
	"net/http"
	"time"
)

var sessions = make(map[string]string)

func login(w http.ResponseWriter) {
	// return login.html page
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t, err := template.ParseFiles("login.html")
	if err != nil {
		fmt.Fprintf(w, "Unable to load template")
	}
	t.Execute(w, nil)
}

func submit(w http.ResponseWriter, r *http.Request) {
	// parse form, then check username and password
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	if username != "admin" || password != "123456" {
		http.Redirect(w, r, "/login/", http.StatusSeeOther)
		return
	}
	// login succeed, set cookie to client and save session
	id := string(uuid.Rand().Hex())
	http.SetCookie(w, &http.Cookie{
		Name:    "SESSIONID",
		Value:   id,
		Path:    "/",
		Expires: time.Now().Add(10 * time.Minute),
		MaxAge:  90000,
	})
	sessions[id] = id
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
