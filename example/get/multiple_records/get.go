package main

import (
	"context"
	"encoding/json"
	"fmt"
	pkg "github.com/megadata-dev/routerosv7-restfull-api"
)

// Create constants for the default values for this example application
const (
	defaultHost     = "10.90.0.251" // Change this to the IP address of your RouterOS device
	defaultUsername = "userapi"     // Change this to the username of your RouterOS device
	defaultPassword = "password"    // Change this to the password of your RouterOS device
	defaultCommand  = "ip/address"  // Change this to the command you want to execute
)

// AppConfig struct for this example application configuration values
type AppConfig struct {
	Host     string // IP address of the RouterOS device
	Username string // Username of the RouterOS device
	Password string // Password of the RouterOS device
	Command  string // Command to execute
}

// NewAppConfig function to create new AppConfig instance with default values for this example application
func NewAppConfig(host, username, password, command string) *AppConfig {
	return &AppConfig{
		Host:     host,
		Username: username,
		Password: password,
		Command:  command,
	}
}

// DataRetriever interface for retrieving data from RouterOS device
type DataRetriever interface {
	GetData(ctx context.Context) (interface{}, error)
}

// RouterOSDataRetriever struct for retrieving data from RouterOS device
type RouterOSDataRetriever struct {
	config *AppConfig
}

// NewRouterOSDataRetriever function to create new RouterOSDataRetriever instance
func NewRouterOSDataRetriever(config *AppConfig) DataRetriever {
	return &RouterOSDataRetriever{config: config}
}

/*
GetData function to retrieve data from RouterOS device using the config values from the RouterOSDataRetriever
instance that called this function and return the data as interface{} and error
*/
func (r *RouterOSDataRetriever) GetData(ctx context.Context) (interface{}, error) {

	// Retrieve data from RouterOS device using the config values from the RouterOSDataRetriever instance
	data, err := pkg.Print(ctx, r.config.Host, r.config.Username, r.config.Password, r.config.Command)

	// Check if there is an error
	if err != nil {
		return nil, err
	}

	// Return the data as interface{} and nil error
	return data, nil
}

// PrintJSON function to print the data as JSON to the console
func PrintJSON(data interface{}) {

	// Check if the data is not []interface{} type
	response, ok := data.([]interface{})

	// Print error message if the data is not []interface{} type and return from this function
	if !ok {
		fmt.Println("Failed to process data")
		return
	}

	// Process and print the data as JSON
	var resultJSON []map[string]interface{}

	// Iterate over the response and append the data to ipAddresses
	for _, item := range response {
		if dataItem, ok := item.(map[string]interface{}); ok {
			resultJSON = append(resultJSON, dataItem)
		}
	}

	// Marshal the resultJSON to JSON string and print it to the console
	jsonData, err := json.Marshal(map[string]interface{}{
		"code":   200,
		"status": "OK",
		"data":   resultJSON,
	})

	// Print error message if there is an error and return from this function
	if err != nil {
		fmt.Println("Failed to marshal JSON:", err)
		return
	}

	// Print the JSON string to the console
	fmt.Println(string(jsonData))
}

// Main function for this example application
func main() {

	// Create new AppConfig instance with default values for this example application
	config := NewAppConfig(defaultHost, defaultUsername, defaultPassword, defaultCommand)

	// Create a PingManager with host configuration for ping and check if the device is available
	pingManager := pkg.NewPing(config.Host)

	// Check if pingManager
	if pingManager == nil {
		fmt.Println("Failed to create PingManager")
		return
	}

	// Run pingManager and check if the device is available
	err := pingManager.CheckAvailableDevice()

	// Check if there is an error
	if err != nil {
		fmt.Println("Device is not available:", err)
	} else {
		fmt.Println("Device is available")
	}

	// Authenticate to RouterOS device and check if the authentication is success or not
	err = pkg.AuthDevice(context.Background(), config.Host, config.Username, config.Password)
	if err != nil {
		fmt.Println("Authentication failed:", err)
	} else {
		fmt.Println("Authentication success")
	}

	// Create new RouterOSDataRetriever instance with the config values from the AppConfig instance
	dataRetriever := NewRouterOSDataRetriever(config)

	// Retrieve data from RouterOS device using the config values from the RouterOSDataRetriever instance
	data, err := dataRetriever.GetData(context.Background())
	if err != nil {
		fmt.Println("Failed to retrieve data:", err)
		return
	}

	// Print the data as JSON to the console
	PrintJSON(data)
}
