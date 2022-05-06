package main

const (
	InvalidAccount      = "The account is invalid"
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
		return false, InvalidAccount
	}
	if at.IsAdmin {
		return true, ""
	}
	if p, exists = GetPermissionsForApp(name, appName); !exists || !p.Enable {
		return false, NoPermission
	}
	for _, path := range p.ExcludedPaths {
		if path == uri {
			return false, NoPermissionForPath
		}
	}
	return true, ""
}
