package main

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
)

const (
	GatewayLoginUri  = "/gateway/login"
	GatewaySubmitUri = "/gateway/submit"
	IdentityKeyName  = "SESSIONID"
)

type handle struct {
	host string
	port string
}

type Setting struct {
	Port     int            `yaml:"port"`
	Auth     bool           `yaml:"auth"`
	Domain   string         `yaml:"domain"`
	Username string         `yaml:"username"`
	Password string         `yaml:"password"`
	Mapping  map[string]int `yaml:"mapping"`
}

type ErrorResponse struct {
	Code            int    // The code of http error
	Reason          string // The reason of http error
	Message         string // The detail message of error
	RedirectToLogin bool   // needs to redirect to login page
}

var conf Setting

func (h *handle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s\n", r.Method, r.URL)

	if r.RequestURI == GatewayLoginUri || r.RequestURI == GatewayLoginUri+"/" {
		Login(w)
		return
	}
	if r.RequestURI == GatewaySubmitUri || r.RequestURI == GatewaySubmitUri+"/" {
		Submit(w, r)
		return
	}
	exists := hasLogin(r)
	if !exists {
		http.Redirect(w, r, GatewayLoginUri, http.StatusSeeOther)
		return
	}
	port, err := getProxyServerPort(r)
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
	proxy.ServeHTTP(w, r)
}

func getProxyServerPort(r *http.Request) (int, error) {
	name := r.Header.Get("APPLICATION_NAME")
	if name == "" {
		return -1, errors.New("THE APPLICATION MAME NOT IN HEADERS")
	}

	var port int
	var ok bool
	if port, ok = conf.Mapping[name]; !ok {
		return -1, errors.New("NO MATCHING APPLICATION")
	}
	return port, nil
}

func hasLogin(r *http.Request) bool {
	cookie, _ := r.Cookie(IdentityKeyName)
	if cookie == nil {
		return false
	}
	_, ok := Sessions[cookie.Value]
	return ok
}

func returnError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t, err := template.ParseFiles("error.html")
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
	log.Printf("Running on http://localhost:%d\n", conf.Port)
	err := http.ListenAndServe(":"+strconv.Itoa(conf.Port), h)
	if err != nil {
		log.Fatalln("ListenAndServe: ", err)
	}
}

func main() {
	buf, err := ioutil.ReadFile("conf.yaml")
	if err != nil {
		log.Fatalln("cannot read conf.yaml file: \n", err)
	}
	conf = Setting{}
	err = yaml.Unmarshal(buf, &conf)
	if err != nil {
		log.Fatalln("cannot load setting: \n", err)
	}

	runServer()
}
