package main

import (
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
	Port               int           `yaml:"port" json:"port"`
	Domain             string        `yaml:"domain" json:"domain"`
	CookieExpiredHours time.Duration `yaml:"cookie-expired-hours" json:"CookieExpiredHours"`
	Debug              bool          `yaml:"debug" json:"debug"`
}

type Account struct {
	Enable   bool   `yaml:"enable" json:"enable"`
	IsAdmin  bool   `yaml:"is-admin" json:"isAdmin"`
	Password string `yaml:"password" json:"password"`
}

type Application struct {
	Enable bool `yaml:"enable" json:"enable"`
	Port   int  `yaml:"port" json:"port"`
}

type Permission struct {
	Enable        bool     `yaml:"enable" json:"enable"`
	IncludedPaths []string `yaml:"included-paths" json:"includedPaths"`
	ExcludedPaths []string `yaml:"excluded-paths" json:"excludedPaths"`
}

type Setting struct {
	Server       Server                           `yaml:"server" json:"server"`
	Accounts     map[string]Account               `yaml:"accounts" json:"accounts"`         // name -> Account
	Applications map[string]Application           `yaml:"applications" json:"applications"` // name -> Application
	Permissions  map[string]map[string]Permission `yaml:"permissions" json:"permissions"`   // username -> {appName -> Permission}
}

const Env = Dev

var Conf Setting

func LoadConfig(envType EnvType) error {
	logger.Info("Start loading config file, envType: ", envType)
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
	Conf = Setting{}
	err = yaml.Unmarshal(buf, &Conf)
	if err != nil {
		return err
	}
	logger.Infof("Read config %s successfully: \n %s", filename, ToJson(Conf))
	// set log level
	if Conf.Server.Debug {
		logger.SetLevel(logger.DebugLevel)
	}
	return nil
}
