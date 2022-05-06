package main

import (
	"encoding/json"
	logger "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
)

const (
	Prod EnvType = "prod"
	Dev  EnvType = "dev"
	Test EnvType = "test"
)

type EnvType string

type Server struct {
	Port               int           `yaml:"port"`
	Domain             string        `yaml:"domain"`
	CookieExpiredHours time.Duration `yaml:"cookie-expired-hours"`
}

type Account struct {
	Enable   bool   `yaml:"enable"`
	IsAdmin  bool   `yaml:"is-admin"`
	Password string `yaml:"password"`
}

type Application struct {
	Enable bool `yaml:"enable"`
	Port   int  `yaml:"port"`
}

type Permission struct {
	Enable        bool     `yaml:"enable"`
	ExcludedPaths []string `yaml:"excluded-paths"`
}

type Setting struct {
	Server       Server                           `yaml:"server"`
	Accounts     map[string]Account               `yaml:"accounts"`     // name -> Account
	Applications map[string]Application           `yaml:"applications"` // name -> Application
	Permissions  map[string]map[string]Permission `yaml:"permissions"`  // username -> {appName -> Permission}
}

var conf Setting

func LoadConfig(envType EnvType) error {
	logger.Info("Start loading config file, envType:", envType)
	var filename string
	if envType == Prod {
		filename = "conf.yaml"
	} else {
		filename = "conf_" + string(envType) + ".yaml"
	}

	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	conf = Setting{}
	err = yaml.Unmarshal(buf, &conf)
	if err != nil {
		return err
	}
	logger.Infof("Read config %s successfully: \n %s", filename, toJson(conf))
	return nil
}

func toJson(setting Setting) string {
	buf, err := json.Marshal(setting)
	if err != nil {
		logger.Error("toJson error", err.Error())
		return ""
	}
	return string(buf)
}
