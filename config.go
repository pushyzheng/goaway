package main

import (
	logger "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
)

var conf Setting

type Setting struct {
	Port               int            `yaml:"port"`
	Auth               bool           `yaml:"auth"`
	Domain             string         `yaml:"domain"`
	Username           string         `yaml:"username"`
	Password           string         `yaml:"password"`
	CookieExpiredHours time.Duration  `yaml:"cookie-expired-hours"`
	Mapping            map[string]int `yaml:"mapping"`
}

func LoadConfig() error {
	buf, err := ioutil.ReadFile("conf.yaml")
	if err != nil {
		return err
	}
	conf = Setting{}
	err = yaml.Unmarshal(buf, &conf)
	if err != nil {
		return err
	}
	logger.Info("Read config successfully: \n", string(buf))
	return nil
}
