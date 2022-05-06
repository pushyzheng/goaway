package main

import (
	logger "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
)

var conf Setting

type Account struct {
	Enable   bool   `yaml:"enable"`
	IsAdmin  bool   `yaml:"is-admin"`
	Password string `yaml:"password"`
}

type Application struct {
	Enable bool `yaml:"enable"`
	Port   int  `yaml:"port"`
}

type Setting struct {
	Port               int                    `yaml:"port"`
	Auth               bool                   `yaml:"auth"`
	Domain             string                 `yaml:"domain"`
	Accounts           map[string]Account     `yaml:"accounts"`             // name -> Account
	Applications       map[string]Application `yaml:"applications"`         // name -> Application
	CookieExpiredHours time.Duration          `yaml:"cookie-expired-hours"` //
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
