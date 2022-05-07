package main

import (
	"encoding/json"
	logger "github.com/sirupsen/logrus"
)

func ToJson(v interface{}) string {
	buf, err := json.Marshal(v)
	if err != nil {
		logger.Error("toJson error", err.Error())
		return ""
	}
	return string(buf)
}

func Contains(slice []string, e string) bool {
	if slice == nil || len(slice) == 0 {
		return false
	}
	for _, each := range slice {
		if each == e {
			return true
		}
	}
	return false
}
