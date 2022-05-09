package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetPermissions(t *testing.T) {
	p1, ok := GetPermissions("admin")
	assert.True(t, ok)
	assert.NotEmpty(t, p1)

	p2, ok := GetPermissions("unknown")
	assert.False(t, ok)
	assert.Empty(t, p2)
}

func TestGetPermissionsForApp(t *testing.T) {
	p1, ok := GetPermissionsForApp("admin", "flask")
	assert.True(t, ok)
	assert.NotEmpty(t, p1)

	p2, ok := GetPermissionsForApp("admin", "unknown")
	assert.False(t, ok)
	assert.False(t, p2.Enable)

	p3, ok := GetPermissionsForApp("unknown", "flask")
	assert.False(t, ok)
	assert.False(t, p3.Enable)
}

func TestHasPermission(t *testing.T) {
	var ok bool

	ok, _ = HasPermission("admin", "flask", "/foo")
	assert.True(t, ok)

	ok, _ = HasPermission("admin", "flask", "/admin")
	assert.True(t, ok) // admin account must is true

	ok, _ = HasPermission("mark", "flask", "/foo")
	assert.True(t, ok)

	ok, _ = HasPermission("mark", "flask", "/foo2")
	assert.False(t, ok)

	ok, _ = HasPermission("mark", "flask", "/public")
	assert.True(t, ok)
}

func TestHasPermission2(t *testing.T) {
	var ok bool
	var cause string

	ok, cause = HasPermission("Michelle", "flask", "/foo")
	assert.False(t, ok)
	assert.Equal(t, InvalidAccount, cause)

	ok, cause = HasPermission("mark", "flask", "/admin")
	assert.False(t, ok)
	assert.Equal(t, NoPermissionForPath, cause)

	ok, cause = HasPermission("mark", "gin", "/foo")
	assert.False(t, ok)
	assert.Equal(t, NoPermission, cause)

	ok, cause = HasPermission("mark", "spring", "/admin")
	assert.False(t, ok)
	assert.Equal(t, InvalidApplication, cause)
}

func init() {
	err := LoadConfig(Test)
	if err != nil {
		panic(err)
	}
}
