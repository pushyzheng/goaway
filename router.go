package main

import (
	"errors"
	logger "github.com/sirupsen/logrus"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
)

type Handler func(resp http.ResponseWriter, req *http.Request)

type Router struct {
	Handler    Handler
	needsLogin bool
}

var routerMapping map[string]Router

func Route(h *handle, resp http.ResponseWriter, req *http.Request) {
	logger.Infof("[%s] %s", req.Method, req.URL)

	reqUrl := req.RequestURI
	if strings.HasPrefix(reqUrl, "/gateway") {
		for path, router := range routerMapping {
			if !(path == reqUrl || path+"/" == reqUrl) {
				continue
			}
			if router.needsLogin && !HasLogin(req) {
				http.Redirect(resp, req, GatewayLoginUri, http.StatusSeeOther)
				return
			}
			router.Handler(resp, req)
			return
		}
		ReturnError(resp, 404, "Not Found")
	} else if !HasLogin(req) {
		http.Redirect(resp, req, GatewayLoginUri, http.StatusSeeOther)
	} else {
		reverseProxy(h, resp, req)
	}
}

func reverseProxy(h *handle, resp http.ResponseWriter, req *http.Request) {
	port, err := getProxyPort(req)
	if err != nil {
		ReturnError(resp, http.StatusBadRequest, err.Error())
		return
	}
	remote, err := url.Parse("http://" + h.host + ":" + strconv.Itoa(port))
	if err != nil {
		ReturnError(resp, http.StatusBadRequest, err.Error())
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ServeHTTP(resp, req)
}

func getProxyPort(r *http.Request) (int, error) {
	name := r.Header.Get("APPLICATION_NAME")
	if len(name) == 0 {
		return -1, errors.New("THE APPLICATION NAME NOT IN HEADERS")
	}
	if app, ok := Conf.Applications[name]; !ok {
		return -1, errors.New("NO MATCHING APPLICATION")
	} else if !app.Enable {
		return -1, errors.New("APPLICATION IS UNAVAILABLE")
	} else {
		return app.Port, nil
	}
}

func init() {
	routerMapping = map[string]Router{
		GatewayLoginUri:         {Handler: Login},
		GatewaySubmitUri:        {Handler: Submit},
		GatewayConfigUri:        {Handler: GetConfig, needsLogin: true},
		GatewayConfigRefreshUri: {Handler: RefreshConfig, needsLogin: true},
	}
}
