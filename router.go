package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	logger "github.com/sirupsen/logrus"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Handler func(resp http.ResponseWriter, req *http.Request)

type Router struct {
	Handler    Handler
	needsLogin bool
}

var (
	routerMapping map[string]Router
	fsMapping     map[string]http.Handler
	prom          = promhttp.Handler()
)

func Route(resp http.ResponseWriter, req *http.Request) {
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
	} else {
		reverseProxy(resp, req)
	}
}

func reverseProxy(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	appName := req.Header.Get("APPLICATION_NAME")
	reqUrl := req.URL.Path
	logger.Debugln("reverse proxy, appName:", appName)

	app, ok := Conf.Applications[appName]
	if !ok || !app.Enable {
		ReturnError(w, http.StatusBadRequest, "The application is invalid")
		return
	}
	if ContainsPath(app.Public, reqUrl) || ContainsPath(app.Public, All) {
		logger.Debugf("'%s' is public url of [%s], skip verification", reqUrl, appName)
	} else {
		// check permission
		user, login := GetUser(req)
		if !login {
			DirectLogin(w, req)
			return
		}
		if ok, cause := HasPermission(user.Username, appName, reqUrl); !ok {
			logger.Warnf("user(%s) don't have permission, app: %s, uri: %s",
				user.Username, appName, req.RequestURI)
			ReturnError(w, http.StatusForbidden, cause)
			return
		}
	}
	reverse(w, req, appName, app)
	reverseCounter.WithLabelValues(appName).Observe(float64(time.Since(start).Milliseconds()))
}

func reverse(w http.ResponseWriter, req *http.Request, appName string, app Application) {
	if app.ServerType == FileServer {
		if fs, ok := fsMapping[appName]; !ok {
			ReturnError(w, http.StatusInternalServerError, "Cannot find file server")
			return
		} else {
			fs.ServeHTTP(w, req)
		}
	} else if app.ServerType == WebServer {
		remote, err := url.Parse("http://" + req.Host + ":" + strconv.Itoa(app.Port))
		if err != nil {
			ReturnError(w, http.StatusBadRequest, err.Error())
			return
		}
		proxy := httputil.NewSingleHostReverseProxy(remote)
		proxy.ServeHTTP(w, req)
	} else {
		logger.Errorln("unsupported server type: ", app.ServerType)
		ReturnError(w, http.StatusInternalServerError, "unknown error")
	}
}

func InitFs() {
	fsMapping = map[string]http.Handler{}
	for name, app := range Conf.Applications {
		if app.ServerType == FileServer {
			if app.Dir == "" {
				panic("The dir cannot be blank: " + name)
			}
			fsMapping[name] = http.FileServer(http.Dir(app.Dir))
		}
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
