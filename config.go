package main

import (
	"fmt"
	"io/ioutil"
	"time"

	logger "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const (
	Prod       EnvType    = "prod"
	Dev        EnvType    = "dev"
	Test       EnvType    = "test"
	FileServer ServerType = "file"
	WebServer  ServerType = "web" // Restful API or html page
)

type EnvType string

type ServerType string

type Server struct {
	Name               string        `yaml:"name" json:"name"`
	Port               int           `yaml:"port" json:"port"`
	Domain             string        `yaml:"domain" json:"domain"`
	CookieExpiredHours time.Duration `yaml:"cookie-expired-hours" json:"cookieExpiredHours"`
	Debug              bool          `yaml:"debug" json:"debug"`
	PrometheusPath     string        `yaml:"prometheus-path" json:"prometheusPath"`
	Statistics         bool          `yaml:"statistics" json:"statistics"`
}

type Account struct {
	Enable   bool   `yaml:"enable" json:"enable"`
	IsAdmin  bool   `yaml:"is-admin" json:"isAdmin"`
	Password string `yaml:"password" json:"password"`
}

type Application struct {
	Enable     bool       `yaml:"enable" json:"enable"`
	ServerType ServerType `yaml:"server-type" json:"serverType"`
	Port       int        `yaml:"port" json:"port"`
	Dir        string     `yaml:"dir" json:"dir"`
	Public     []string   `yaml:"public" json:"public"`
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

var Env EnvType

var Conf Setting

func LoadConfig(envType EnvType) error {
	Env = envType
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
	if err = check(); err != nil {
		return fmt.Errorf("check fail: " + err.Error())
	}
	// set log level
	if Conf.Server.Debug {
		logger.SetLevel(logger.DebugLevel)
	}
	return nil
}

func check() error {
	var err error
	if err = assertNotBlank("example-server.domain", Conf.Server.Domain); err != nil {
		return err
	}
	return nil
}

func assertNotBlank(name string, s string) error {
	if len(s) == 0 {
		return fmt.Errorf("the %s cannot be blank", name)
	}
	return nil
}
