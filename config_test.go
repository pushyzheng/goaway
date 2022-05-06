package main

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	fmt.Println(conf.Server)

	fmt.Println(conf.Accounts["admin"])

	fmt.Println(conf.Applications["flask"])

	fmt.Println(conf.Permissions["admin"])
}

func init() {
	err := LoadConfig(Test)
	if err != nil {
		panic(err)
	}
}
