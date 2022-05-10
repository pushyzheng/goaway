package main

import (
	logger "github.com/sirupsen/logrus"
	"log"
	"net/http"
	"strconv"
	"time"
)

type handle struct {
	host string
	port string
}

func (h *handle) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	start := time.Now()
	Route(resp, req)
	cost := time.Since(start)
	entranceSummary.Observe(float64(cost.Milliseconds()))
	logger.Infof("[Proxy] | %dms | %s %s", cost.Milliseconds(), req.Method, req.URL)
}

func runServer() {
	h := &handle{}
	logger.Infof("Running on http://localhost:%d", Conf.Server.Port)
	err := http.ListenAndServe(":"+strconv.Itoa(Conf.Server.Port), h)
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
