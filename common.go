package main

import (
	"encoding/json"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"html/template"
	"log"
	"net/http"
)

const (
	GatewayUriPrefix        = "/gateway"
	GatewayLoginUri         = GatewayUriPrefix + "/login"
	GatewaySubmitUri        = GatewayUriPrefix + "/submit"
	GatewayConfigUri        = GatewayUriPrefix + "/config"
	GatewayConfigRefreshUri = GatewayUriPrefix + "/config/refresh"
	IdentityKeyName         = "SESSION_ID"
)

type ErrorResponse struct {
	Code            int    // The code of http error
	Reason          string // The reason of http error
	Message         string // The detail message of error
	RedirectToLogin bool   // needs to redirect to login page
}

type APIResponse struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func ReturnJson(resp http.ResponseWriter, data interface{}) {
	ar := APIResponse{Code: 0, Data: data}
	buf, err := json.Marshal(ar)
	if err != nil {
		logger.Errorln("json marshal error:", ar)
		buf = []byte("json marshal error")
	}
	_, err = resp.Write(buf)
	if err != nil {
		logger.Errorln("write data error:", string(buf))
	}
}

func ReturnError(w http.ResponseWriter, code int, msg string) {
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
	_ = t.Execute(w, resp)
}
