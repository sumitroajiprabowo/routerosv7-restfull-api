package main

import (
	"fmt"
	pkg "github.com/megadata-dev/routerosv7-restfull-api"
)

const (
	hostToPing = "192.168.88.1" // Change with your host
)

func main() {
	pingDevice(hostToPing)
}

func pingDevice(host string) {
	// Create a PingManager with host configuration for ping
	pingManager := pkg.NewPing(host)

	// Check if pingManager is nil
	if pingManager == nil {
		fmt.Println("Failed to create PingManager") // Print error message
		return
	}

	// Run pingManager and check if the device is available
	err := pingManager.CheckAvailableDevice()

	if err != nil {
		fmt.Println("Device is not available:", err)
	} else {
		fmt.Println("Device is available")
	}
}
