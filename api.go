package main

import (
	logger "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func GetConfig(resp http.ResponseWriter, req *http.Request) {
	ReturnJson(resp, Conf)
}

func GetSessions(resp http.ResponseWriter, req *http.Request) {
	tmpMap := make(map[string]User)
	sessions.Range(func(k, v interface{}) bool {
		tmpMap[k.(string)] = v.(User)
		return true
	})
	ReturnJson(resp, tmpMap)
}

func RefreshConfig(resp http.ResponseWriter, req *http.Request) {
	err := LoadConfig(Env)
	if err != nil {
		logger.Errorln("refresh config error:", err)
		return
	}
	ReturnJson(resp, "success")
}

func GetStatisticsPage(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "text/html; charset=utf-8")
	buf, err := ioutil.ReadFile(StatisticsPagePath)
	if err != nil {
		logger.Errorf("load page error, path = %s, err = %s", StatisticsPagePath, err.Error())
		buf = []byte("json marshal error")
	}
	_, _ = resp.Write(buf)
}

func GetStatistics(resp http.ResponseWriter, req *http.Request) {
	apps, err := GetApps()
	res := make(map[string]interface{})
	if err != nil {
		ReturnError(resp, 503, "Cannot get apps")
		return
	}
	logger.Info(apps)

	// set refers
	refersMapping := make(map[string]interface{})
	for _, appName := range apps {
		refers, err := GetRefers(appName)
		if err != nil {
			logger.Errorf("cannot get refers, name: %s", appName)
		}
		refersMapping[appName] = refers
	}
	// setPaths
	pathMapping := make(map[string]interface{})
	for _, appName := range apps {
		refers, err := GetPaths(appName)
		if err != nil {
			logger.Errorf("cannot get paths, name: %s", appName)
		}
		pathMapping[appName] = refers
	}

	res["refers"] = refersMapping
	res["paths"] = pathMapping
	ReturnJson(resp, res)
}
