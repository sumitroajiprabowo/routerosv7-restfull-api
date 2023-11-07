package main

import (
	"context"
	"encoding/json"
	"fmt"
	pkg "github.com/megadata-dev/routerosv7-restfull-api"
)

func getData() []interface{} {
	cmd := fmt.Sprintf("ip/address")

	// Get all IP address
	data, err := pkg.Print(context.Background(), "192.168.88.1", "userapi", "password", cmd)

	if err != nil {
		return nil
	}

	// Convert interface to []interface{}
	response := data.([]interface{})

	// Create a slice of maps to hold the structured data
	var ipAddresses []map[string]interface{}

	// Convert the data to the desired format
	for _, item := range response {
		if ipItem, ok := item.(map[string]interface{}); ok {
			ipAddresses = append(ipAddresses, ipItem)
		}
	}

	// Print the JSON representation
	jsonData, _ := json.Marshal(map[string]interface{}{
		"code":   200,
		"status": "OK",
		"data":   ipAddresses,
	})

	fmt.Println(string(jsonData))

	return response
}

func main() {

	// Check Available Device
	err := pkg.PingDevice("192.168.88.1")
	if err != nil {
		fmt.Println("Device is not available:", err)
	} else {
		fmt.Println("Device is available")
	}

	// Authentication
	err = pkg.AuthDevice(context.Background(), "192.168.88.1", "userapi", "password")
	if err != nil {
		fmt.Println("Authentication failed:", err)
	} else {
		fmt.Println("Authentication success")
	}

	// Get Data
	getData()

	return

}
