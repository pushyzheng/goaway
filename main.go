package main

import (
	"flag"
	logger "github.com/sirupsen/logrus"
	"log"
	"net/http"
	"strconv"
	"time"
)

var envFlag = flag.String("env", string(Dev), "Input env type")

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
	flag.Parse()
	err := LoadConfig(parseEnv(*envFlag))
	if err != nil {
		panic("cannot load config: " + err.Error())
	}
	InitFs()
	RegisterProm()
	initRedis()
	runServer()
}

func parseEnv(s string) EnvType {
	if s == string(Prod) {
		return Prod
	} else if s == string(Dev) {
		return Dev
	} else if s == string(Test) {
		return Test
	} else {
		panic("The env type is illegal: " + s)
	}
}

func init() {
	logger.SetFormatter(&logger.TextFormatter{
		FullTimestamp: true,
	})
}
