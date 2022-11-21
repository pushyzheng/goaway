package main

import (
	"fmt"
	"net/http"
)

const AppsKey = "goaway::statistics::apps"

func Record(appName string, req *http.Request) {
	refer := req.Header.Get("Referer")
	uri := req.URL.Path

	incrOne("goaway::statistics::apps", appName)
	if len(refer) > 0 {
		incrOne(formatReferKey(appName), ParseDomainFromUrl(refer))
	}
	incrOne(formatPathKey(appName), uri)
	incrOne("goaway::statistics::date", GetTodayDate())
}

func GetApps() ([]string, error) {
	return redisConn.HKeys(AppsKey).Result()
}

func GetRefers(appName string) (map[string]string, error) {
	return redisConn.HGetAll(formatReferKey(appName)).Result()
}

func GetPaths(appName string) (map[string]string, error) {
	return redisConn.HGetAll(formatPathKey(appName)).Result()
}

func formatReferKey(appName string) string {
	return fmt.Sprintf("goaway::statistics::entry_refer_%s", appName)
}

func formatPathKey(appName string) string {
	return fmt.Sprintf("goaway::statistics::entry_path_%s", appName)
}

func incrOne(key string, field string) {
	redisConn.HIncrBy(key, field, 1)
}
