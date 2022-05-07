package main

import (
	logger "github.com/sirupsen/logrus"
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
