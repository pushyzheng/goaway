package main

import (
	"encoding/json"
	logger "github.com/sirupsen/logrus"
	"net/url"
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
	return ContainsPredicate(slice, e, func(each, other string) bool {
		return each == other
	})
}

func ContainsPath(slice []string, uri string) bool {
	return ContainsPredicate(slice, uri, func(each, other string) bool {
		ee := EncodeUrl(each)
		return ee == other || ee == other+"/" || ee+"/" == other
	})
}

func ContainsPredicate(slice []string, e string, predicate func(each, other string) bool) bool {
	if slice == nil || len(slice) == 0 {
		return false
	}
	for _, each := range slice {
		if predicate(each, e) {
			return true
		}
	}
	return false
}

func EncodeUrl(s string) string {
	u, err := url.Parse(s)
	if err != nil {
		logger.Error("EncodeUrl error: ", s)
		return s
	}
	return u.EscapedPath()
}

func EncodeUrlComponent(s string) string {
	return url.QueryEscape(s)
}

func EqualsUri(uri, other string) bool {
	return uri == other || uri == other+"/"
}
