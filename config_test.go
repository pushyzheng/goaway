package main

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	fmt.Println(conf.Accounts["admin"])
	fmt.Println(conf.Applications["flask"])
}

func init() {
	err := LoadConfig()
	if err != nil {
		panic(err)
	}
}
