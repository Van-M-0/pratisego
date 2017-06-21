package main

import (
	"gateway"
	"login"
	"sync"
	"os"
	"fmt"
)

func main() {

	serverType := os.Args[1]
	if serverType == "gate" {
		fmt.Print("gate start")
		gw := gateway.NewGateway()
		gw.Start()
	} else if serverType == "login" {
		fmt.Println("login start")
		ls := login.NewLoginServer()
		ls.Start()
	} else if serverType == "client" {

	} else if serverType == "mysql" {
		go mysql_test()
	} else {
		fmt.Println("aaaaaaaa")
	}

	if serverType != "client" {
		wg := new(sync.WaitGroup)
		wg.Add(1)
		wg.Wait()
	}
}
