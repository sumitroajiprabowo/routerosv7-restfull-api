package main

import (
	"context"
	"encoding/json"
	"fmt"
	pkg "github.com/megadata-dev/routerosv7-restfull-api"
	"strings"
)

// Create constants for the default host, username, and password.
const (
	defaultHost     = "10.90.0.251" // defaultHost is the default host for the RouterOS device
	defaultUsername = "userapi"     // defaultUsername is the default username for the RouterOS device
	defaultPassword = "password"    // defaultPassword is the default password for the RouterOS device
)

// AppConfig is a struct that contains the configuration for the application.
type AppConfig struct {
	Host     string            // Host is the IP address of the RouterOS device
	Username string            // Username is the username for the RouterOS device
	Password string            // Password is the password for the RouterOS device
	Command  string            // Command is the command to be executed
	Params   map[string]string //Params is the parameter for the command
}

/*
NewAppConfig is a function to create a new AppConfig.
It accepts the host, username, password, and command as parameters.
It returns a pointer to an AppConfig.
*/
func NewAppConfig(host, username, password, command string, params map[string]string) *AppConfig {
	return &AppConfig{
		Host:     host,     // Set the host
		Username: username, // Set the username
		Password: password, // Set the password
		Command:  command,  // Set the command
		Params:   params,   // Set the params
	}
}

/*
DataRetriever is an interface that contains the GetData method.
This method will be used to retrieve data from the RouterOS device.
The GetData method accepts a context.Context as a parameter and returns an interface and an error.
*/
type DataRetriever interface {
	GetData(ctx context.Context) (interface{}, error)
}

// RouterOSDataRetriever is a struct that implements the DataRetriever interface.
type RouterOSDataRetriever struct {
	config *AppConfig // config is a pointer to an AppConfig
}

/*
NewRouterOSDataRetriever is a function to create a new RouterOSDataRetriever.
It accepts an AppConfig as a parameter and returns a DataRetriever.
The function returns a pointer to a RouterOSDataRetriever.
*/
func NewRouterOSDataRetriever(config *AppConfig) DataRetriever {
	return &RouterOSDataRetriever{config: config}
}

/*
GetData is a method to retrieve data from the RouterOS device.
for Example i want request with single record with multiple params in Mikrotik RouterOS
you can use this function to get data from Mikrotik RouterOS
*/
func (r *RouterOSDataRetriever) GetData(ctx context.Context) (interface{}, error) {

	// Create a slice of strings to store the query parameters
	var queryParams []string

	// Loop through the params and append them to the slice of strings as key=value
	for key, value := range r.config.Params {
		queryParams = append(queryParams, fmt.Sprintf("%s=%s", key, value)) // Set the query parameters
	}

	// Create the command by joining the query parameters with & and adding them to the command string
	command := fmt.Sprintf(r.config.Command, strings.Join(queryParams, "&"))

	// Retrieve the data from the RouterOS device using the Print function from the routerosv7_restfull_api package
	data, err := pkg.Print(ctx, r.config.Host, r.config.Username, r.config.Password, command)
	if err != nil {
		return nil, err
	}

	// Return the data and nil error
	return data, nil
}

/*
PrintJSON is a function to print the data as JSON.
in this function i use interface{}
After that, I loop through the response and convert each item to a map[string]interface{}.
Finally, I marshal the data to JSON and print it to the console.
*/
func PrintJSON(data interface{}) {
	// Convert the data to a slice of interfaces
	response, ok := data.([]interface{})

	// Check if the conversion was successful
	if !ok {
		fmt.Println("Failed to process data")
		return
	}

	// Process and print the data as JSON by converting each item to a map[string]interface{}
	var ipAddresses []map[string]interface{}
	for _, item := range response {
		if ipItem, ok := item.(map[string]interface{}); ok {
			ipAddresses = append(ipAddresses, ipItem)
		}
	}

	// Marshal the data to JSON
	jsonData, err := json.Marshal(map[string]interface{}{
		"code":   200,
		"status": "OK",
		"data":   ipAddresses,
	})
	if err != nil {
		fmt.Println("Failed to marshal JSON:", err)
		return
	}

	// Print the JSON
	fmt.Println(string(jsonData))
}

func main() {

	// Set the parameters for the command to be executed
	params := map[string]string{
		"address":  "192.168.88.1/24",
		"disabled": "false",
	}

	/*
		Create a new AppConfig with the default host, username, password, command, and params
		In this example, the command is ip/address and the params are address=192.168.88.1/24 and disabled=false
		if i use curl command, the command will be like this:
		curl -k -u username:password "http://x.x.x.x/rest/ip/address?address=192.168.88.1/24&disabled=false" |jq
	*/
	config := NewAppConfig(defaultHost, defaultUsername, defaultPassword, "ip/address?%s", params)

	// Create a new RouterOSDataRetriever with the config
	pingManager := pkg.NewPing(config.Host)

	// Check if the device is available
	if pingManager == nil {
		fmt.Println("Failed to create PingManager")
		return
	}

	// Check if the device is available by calling the CheckAvailableDevice method
	err := pingManager.CheckAvailableDevice()
	if err != nil {
		fmt.Println("Device is not available:", err)
		return
	}

	fmt.Println("Device is available")

	// Authenticate to the RouterOS device by calling the AuthDevice method
	err = pkg.AuthDevice(context.Background(), config.Host, config.Username, config.Password)
	if err != nil {
		fmt.Println("Authentication failed:", err)
		return
	}

	fmt.Println("Authentication success")

	// Create a new RouterOSDataRetriever with the config and call the GetData method to retrieve the data
	dataRetriever := NewRouterOSDataRetriever(config)

	// Retrieve the data from the RouterOS device by calling the GetData method
	data, err := dataRetriever.GetData(context.Background())
	if err != nil {
		fmt.Println("Failed to retrieve data:", err)
		return
	}

	PrintJSON(data)
}
