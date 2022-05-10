package main

import (
	"errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

var (
	routerMapping map[string]Router
	prom          = promhttp.Handler()
)

func Route(h *handle, resp http.ResponseWriter, req *http.Request) {
	reqUrl := req.URL.Path
	routerCounter.WithLabelValues(reqUrl).Inc()

	if Conf.Server.PrometheusPath != "" && EqualsUri(reqUrl, Conf.Server.PrometheusPath) {
		prom.ServeHTTP(resp, req)
		return
	}
	if strings.HasPrefix(reqUrl, GatewayUriPrefix) {
		for path, router := range routerMapping {
			if !EqualsUri(reqUrl, path) {
				continue
			}
			if router.needsLogin && !HasLogin(req) {
				DirectLogin(resp, req)
				return
			}
			router.Handler(resp, req)
			return
		}
		ReturnError(resp, 404, "Not Found")
	} else if !HasLogin(req) {
		logger.Debugln("redirect to login page, reqUrl:", reqUrl)
		DirectLogin(resp, req)
	} else {
		reverseProxy(h, resp, req)
	}
}

func reverseProxy(h *handle, resp http.ResponseWriter, req *http.Request) {
	appName := req.Header.Get("APPLICATION_NAME")
	logger.Debugln("reverse proxy, appName:", appName)

	if len(appName) == 0 {
		reverseCounter.WithLabelValues("null").Inc()
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
		ReturnError(resp, http.StatusForbidden, cause)
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
	reverseCounter.WithLabelValues(appName).Inc()
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ServeHTTP(resp, req)
}

func getProxyPort(name string) (int, error) {
	if app, ok := Conf.Applications[name]; !ok {
		return -1, errors.New("no matching application")
	} else if !app.Enable {
		return -1, errors.New("application is unavailable")
	} else {
		return app.Port, nil
	}
}

func init() {
	routerMapping = map[string]Router{
		GatewayLoginUri:         {Handler: Login},
		GatewaySubmitUri:        {Handler: Submit},
		GatewayLogoutUri:        {Handler: Logout},
		GatewayConfigUri:        {Handler: GetConfig, needsLogin: true},
		GatewayConfigRefreshUri: {Handler: RefreshConfig, needsLogin: true},
		GatewaySessionUri:       {Handler: GetSessions, needsLogin: true},
	}
}
