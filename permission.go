package main

import (
	logger "github.com/sirupsen/logrus"
)

const (
	All                 = "*"
	InvalidAccount      = "The account is invalid"
	InvalidApplication  = "The application is invalid"
	NoPermission        = "No permission for this application"
	NoPermissionForPath = "No permission for this path"
)

var empty = Permission{}

func GetPermissions(name string) (map[string]Permission, bool) {
	if len(name) == 0 {
		return nil, false
	}
	if p, ok := Conf.Permissions[name]; !ok {
		return nil, false
	} else {
		return p, true
	}
}

func GetPermissionsForApp(name string, appName string) (Permission, bool) {
	if len(name) == 0 || len(appName) == 0 {
		return empty, false
	}
	if ps, ok := GetPermissions(name); !ok {
		return empty, false
	} else if pa, ok := ps[appName]; !ok {
		return empty, false
	} else {
		return pa, true
	}
}

func HasPermission(name string, appName string, uri string) (bool, string) {
	var at Account
	var p Permission
	var exists bool

	if at, exists = Conf.Accounts[name]; !exists || !at.Enable {
		logger.Warn("The account is invalid: ", appName)
		return false, InvalidAccount
	}
	if at.IsAdmin {
		logger.Debugf("The user(%s) is admin, skip permission verification", name)
		return true, ""
	}
	// check the public path of application
	if app, exists := Conf.Applications[appName]; !exists || !app.Enable {
		logger.Warn("The application is invalid: ", appName)
		return false, InvalidApplication
	} else if ContainsPath(app.Public, All) || ContainsPath(app.Public, uri) {
		return true, ""
	}

	// check the permission of current user
	if p, exists = GetPermissionsForApp(name, appName); !exists || !p.Enable {
		return false, NoPermission
	}
	logger.Debugf("The permissions of current user(%s):%s", name, ToJson(p))

	if (ContainsPath(p.IncludedPaths, All) || ContainsPath(p.IncludedPaths, uri)) &&
		(!ContainsPath(p.ExcludedPaths, All) || !ContainsPath(p.ExcludedPaths, uri)) {
		return true, ""
	}
	return false, NoPermissionForPath
}
