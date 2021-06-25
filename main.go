package main

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
)

type handle struct {
	host string
	port string
}

type Setting struct {
	Port    int            `yaml:"port"`
	auth    bool           `yaml:"auth"`
	Mapping map[string]int `yaml:"mapping"`
}

var conf Setting

func (h *handle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s\n", r.Method, r.URL)

	if r.RequestURI == "/login" || r.RequestURI == "/login/" {
		Login(w)
		return
	}
	if r.RequestURI == "/submit" || r.RequestURI == "/submit/" {
		Submit(w, r)
		return
	}
	cookie, _ := r.Cookie("SESSIONID")
	if cookie != nil {
		sessionId := cookie.Value
		if _, ok := Sessions[sessionId]; ok {
			port, err := getProxyServerPort(r)
			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}
			remote, err := url.Parse("http://" + h.host + ":" + strconv.Itoa(port))
			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}
			proxy := httputil.NewSingleHostReverseProxy(remote)
			proxy.ServeHTTP(w, r)
			return
		}
	}
	http.Redirect(w, r, "/login/", http.StatusSeeOther)
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
