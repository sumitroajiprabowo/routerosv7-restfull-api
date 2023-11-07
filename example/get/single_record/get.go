package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/megadata-dev/routerosv7-restfull-api"
	"strings"
)

// Create constants for the default values for this example application
const (
	routerIP = "10.90.0.251" // Change this to your router's IP address
	username = "userapi"     // Change this to your router's username
	password = "password"    // Change this to your router's password
)

// AppConfig struct for this example application configuration values
type AppConfig struct {
	Host     string            // IP address of the RouterOS device
	Username string            // Username of the RouterOS device
	Password string            // Password of the RouterOS device
	Command  string            // Command to execute
	Params   map[string]string // Params for the command
}

// NewAppConfig function to create new AppConfig instance with default values for this example application
func NewAppConfig(host, username, password, command string, params map[string]string) *AppConfig {
	return &AppConfig{
		Host:     host,     // IP address of the RouterOS device
		Username: username, // Username of the RouterOS device
		Password: password, // Password of the RouterOS device
		Command:  command,  // Command to execute
		Params:   params,   // Params for the command
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

// authenticate authenticates to the router. If authentication fails, an error is returned.
func authenticate(routerIP, username, password string) error {
	// Create a config for the router to authenticate
	config := routerosv7_restfull_api.AuthDeviceConfig{
		Host:     routerIP, // Change this to your router's IP address
		Username: username, // Change this to your router's username
		Password: password, // Change this to your router's password
	}

	// Authenticate to the router
	_, err := routerosv7_restfull_api.AuthDevice(context.Background(), config)
	return err
}

// GetData function to retrieve data from RouterOS device using the config values from the RouterOSDataRetriever
func (r *RouterOSDataRetriever) GetData(ctx context.Context) (interface{}, error) {

	// Create queryParams variable as slice of string
	var queryParams []string

	// Loop through the params and append the key and value to the queryParams slice
	for key, value := range r.config.Params {
		queryParams = append(queryParams, fmt.Sprintf("%s=%s", key, value))
	}

	// Create command variable with the command and queryParams
	command := fmt.Sprintf(r.config.Command, strings.Join(queryParams, "&"))

	// Retrieve data from RouterOS device using the config values from the RouterOSDataRetriever instance
	data, err := routerosv7_restfull_api.Print(ctx, r.config.Host, r.config.Username, r.config.Password, command)

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

	// Iterate over the response and append the data to resultJSON
	for _, item := range response {
		if ipItem, ok := item.(map[string]interface{}); ok {
			resultJSON = append(resultJSON, ipItem)
		}
	}

	// Marshal the resultJSON to JSON
	jsonData, err := json.Marshal(map[string]interface{}{
		"code":   200,        // HTTP status code
		"status": "OK",       // HTTP status message
		"data":   resultJSON, // Data
	})

	// Print error message if there is an error and return from this function
	if err != nil {
		fmt.Println("Failed to marshal JSON:", err)
		return
	}

	// Print the JSON string to the console
	fmt.Println(string(jsonData))
}

// main function for this example application
func main() {

	// Create params variable as map[string]string
	params := map[string]string{
		"address":  "192.168.88.1/24",
		"disabled": "false",
	}

	// Create config variable with the default values and params
	config := NewAppConfig(routerIP, username, password, "ip/address?%s", params)

	// Create a PingManager with host configuration for ping
	pingManager := routerosv7_restfull_api.NewPing(config.Host)

	// Check if pingManager is nil
	if pingManager == nil {
		fmt.Println("Failed to create PingManager")
		return
	}

	// Run pingManager and check if the device is available
	err := pingManager.CheckAvailableDevice()

	// Check if there is an error
	if err != nil {
		fmt.Println("Device is not available:", err)
		return
	}

	fmt.Println("Device is available")

	// Authenticate to the router using the config values from the AppConfig instance
	err = authenticate(routerIP, username, password)

	// Check if authentication failed
	if err != nil {
		fmt.Println("Authentication failed:", err)
		return
	}

	fmt.Println("Authentication success")

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
