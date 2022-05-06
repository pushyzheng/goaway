package main

import (
	logger "github.com/sirupsen/logrus"
	"log"
	"net/http"
	"strconv"
)

type handle struct {
	host string
	port string
}

func (h *handle) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	Route(h, resp, req)
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
	err := LoadConfig(Prod)
	if err != nil {
		panic("cannot load config: " + err.Error())
	}
	runServer()
}

func init() {
	logger.SetFormatter(&logger.TextFormatter{
		FullTimestamp: true,
	})
}
