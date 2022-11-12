package main

import (
	"fmt"
	"net/http"
)

func record(appName string, req *http.Request) {
	refer := req.Header.Get("Referer")
	uri := req.URL.Path

	redisConn.HIncrBy("goaway::statistics::apps", appName, 1)
	if len(refer) > 0 {
		hostname := ParseDomainFromUrl(refer)
		redisConn.HIncrBy(fmt.Sprintf("goaway::statistics::entry_refer_%s", appName), hostname, 1)
	}
	redisConn.HIncrBy(fmt.Sprintf("goaway::statistics::entry_path_%s", appName), uri, 1)
}
