package main

import (
	"context"
	"fmt"
	"github.com/megadata-dev/routerosv7-restfull-api"
)

// Create a constant for the router's IP address, username, and password
const (
	routerIP = "10.90.0.251" // Change this to your router's IP address
	username = "userapi"     // Change this to your router's username
	password = "password"    // Change this to your router's password
)

// The main function
func main() {

	// Authenticate to the router
	err := authenticate(routerIP, username, password)

	// Check if authentication failed
	if err != nil {
		fmt.Println("Authentication failed:", err)
	} else {
		fmt.Println("Authentication success")
	}
}

// authenticate authenticates to the router. If authentication fails, an error is returned.
func authenticate(routerIP, username, password string) error {

	// Create a config for the router to authenticate
	config := routerosv7_restfull_api.AuthConfig{
		Host:     routerIP, // Change this to your router's IP address
		Username: username, // Change this to your router's username
		Password: password, // Change this to your router's password
	}

	// Authenticate to the router
	_, err := routerosv7_restfull_api.Auth(context.Background(), config)
	return err
}
