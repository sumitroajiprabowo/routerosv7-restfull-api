package main

import (
	"context"
	"fmt"
	pkg "github.com/megadata-dev/routerosv7-restfull-api"
)

const (
	routerIP = "10.90.0.251"
	username = "userapi"
	password = "password"
)

func main() {
	err := authenticate(routerIP, username, password)
	if err != nil {
		fmt.Println("Authentication failed:", err)
	} else {
		fmt.Println("Authentication success")
	}
}

func authenticate(routerIP, username, password string) error {
	return pkg.AuthDevice(context.Background(), routerIP, username, password)
}
