package main

import (
	"errors"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"html/template"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
)

const (
	GatewayLoginUri  = "/gateway/login"
	GatewaySubmitUri = "/gateway/submit"
	IdentityKeyName  = "SESSION_ID"
)

type handle struct {
	host string
	port string
}

type ErrorResponse struct {
	Code            int    // The code of http error
	Reason          string // The reason of http error
	Message         string // The detail message of error
	RedirectToLogin bool   // needs to redirect to login page
}

func (h *handle) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	logger.Infof("[%s] %s", req.Method, req.URL)

	if req.RequestURI == GatewayLoginUri || req.RequestURI == GatewayLoginUri+"/" {
		Login(w)
		return
	}
	if req.RequestURI == GatewaySubmitUri || req.RequestURI == GatewaySubmitUri+"/" {
		Submit(w, req)
		return
	}
	if !HasLogin(req) {
		http.Redirect(w, req, GatewayLoginUri, http.StatusSeeOther)
		return
	}
	port, err := getProxyPort(req)
	if err != nil {
		returnError(w, http.StatusBadRequest, err.Error())
		return
	}
	remote, err := url.Parse("http://" + h.host + ":" + strconv.Itoa(port))
	if err != nil {
		returnError(w, http.StatusBadRequest, err.Error())
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ServeHTTP(w, req)
}

func getProxyPort(r *http.Request) (int, error) {
	name := r.Header.Get("APPLICATION_NAME")
	if len(name) == 0 {
		return -1, errors.New("THE APPLICATION NAME NOT IN HEADERS")
	}
	if app, ok := conf.Applications[name]; !ok {
		return -1, errors.New("NO MATCHING APPLICATION")
	} else if !app.Enable {
		return -1, errors.New("APPLICATION IS UNAVAILABLE")
	} else {
		return app.Port, nil
	}
}

func returnError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t, err := template.ParseFiles("static/error.html")
	if err != nil {
		log.Printf("parse file error: %s\n", err.Error())
		_, _ = fmt.Fprintf(w, "Unable to load template")
		return
	}
	resp := ErrorResponse{
		Code:    code,
		Message: msg,
		Reason:  http.StatusText(code),
	}
	if code == http.StatusUnauthorized {
		resp.RedirectToLogin = true
	}
	t.Execute(w, resp)
}

func runServer() {
	h := &handle{}
	logger.Infof("Running on http://localhost:%d", conf.Port)
	err := http.ListenAndServe(":"+strconv.Itoa(conf.Port), h)
	if err != nil {
		log.Fatalln("ListenAndServe: ", err)
	}
}

func main() {
	logger.SetFormatter(&logger.TextFormatter{
		FullTimestamp: true,
	})

	err := LoadConfig()
	if err != nil {
		panic("cannot load config: " + err.Error())
	}
	runServer()
}
