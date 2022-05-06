package main

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	fmt.Println(Conf.Server)

	fmt.Println(Conf.Accounts["admin"])

	fmt.Println(Conf.Applications["flask"])

	fmt.Println(Conf.Permissions["admin"])
}

func init() {
	err := LoadConfig(Test)
	if err != nil {
		panic(err)
	}
}
