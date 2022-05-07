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
	if strings.HasPrefix(reqUrl, GatewayUriPrefix) {
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
		logger.Debugln("redirect to login page, reqUrl:", reqUrl)
		http.Redirect(resp, req, GatewayLoginUri, http.StatusSeeOther)
	} else {
		reverseProxy(h, resp, req)
	}
}

func reverseProxy(h *handle, resp http.ResponseWriter, req *http.Request) {
	appName := req.Header.Get("APPLICATION_NAME")
	logger.Debugln("reverse proxy, appName:", appName)

	if len(appName) == 0 {
		ReturnError(resp, http.StatusBadRequest, "Cannot get the name of application in headers")
		return
	}
	// check permission
	user, login := GetUser(req)
	if !login {
		ReturnError(resp, http.StatusUnauthorized, "The user don't login")
		return
	}
	if ok, cause := HasPermission(user.Username, appName, req.RequestURI); !ok {
		logger.Warnf("user(%s) don't have permission, app: %s, uri: %s",
			user.Username, appName, req.RequestURI)
		ReturnError(resp, http.StatusUnauthorized, cause)
		return
	}
	port, err := getProxyPort(appName)
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

func getProxyPort(name string) (int, error) {
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
		GatewaySessionUri:       {Handler: GetSessions, needsLogin: true},
	}
}
