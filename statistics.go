package main

import (
	"fmt"
	"net/http"
)

func record(appName string, req *http.Request) {
	referArr := req.Header["refer"]
	uri := req.URL.Path

	redisConn.HIncrBy("goaway::statistics::apps", appName, 1)
	if len(referArr) != 0 {
		refer := referArr[0]
		redisConn.HIncrBy(fmt.Sprintf("goaway::statistics::entry_refer_%s", appName), refer, 1)
	}
	redisConn.HIncrBy(fmt.Sprintf("goaway::statistics::entry_path_%s", appName), uri, 1)
}
