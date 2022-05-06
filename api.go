package main

import (
	logger "github.com/sirupsen/logrus"
	"net/http"
)

func GetConfig(resp http.ResponseWriter, req *http.Request) {
	ReturnJson(resp, Conf)
}

func RefreshConfig(resp http.ResponseWriter, req *http.Request) {
	err := LoadConfig(Prod)
	if err != nil {
		logger.Errorln("refresh config error:", err)
		return
	}
	ReturnJson(resp, "success")
}
