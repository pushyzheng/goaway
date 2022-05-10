package main

import (
	logger "github.com/sirupsen/logrus"
	"log"
	"net/http"
	"strconv"
	"time"
)

var sh ServerHandler

type ServerHandler struct {
	Host string
	Port string
}

func (h *ServerHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	start := time.Now()
	Route(resp, req)
	cost := time.Since(start)
	entranceSummary.Observe(float64(cost.Milliseconds()))
	logger.Infof("[Proxy] | %dms | %s %s", cost.Milliseconds(), req.Method, req.URL)
}

func runServer() {
	sh = ServerHandler{}
	logger.Infof("Running on http://localhost:%d", Conf.Server.Port)
	err := http.ListenAndServe(":"+strconv.Itoa(Conf.Server.Port), &sh)
	if err != nil {
		log.Fatalln("ListenAndServe: ", err)
	}
}

func main() {
	err := LoadConfig(Env)
	if err != nil {
		panic("cannot load config: " + err.Error())
	}
	InitFs()
	runServer()
}

func init() {
	logger.SetFormatter(&logger.TextFormatter{
		FullTimestamp: true,
	})
}
