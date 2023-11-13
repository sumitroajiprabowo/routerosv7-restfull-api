// api_methods.go

package routerosv7_restfull_api

import (
	"context"
	"fmt"
)

// AuthConfig is the configuration for the AuthDevice function
type AuthConfig struct {
	Host     string // Host for the request to Mikrotik Router
	Username string // Username for the request to Mikrotik Router
	Password string // Password for the request to Mikrotik Router
}

// APIRequest represents a request to the API.
type APIRequest struct {
	Host     string // Host for the request to Mikrotik Router
	Username string // Username for the request to Mikrotik Router
	Password string // Password for the request to Mikrotik Router
	Command  string // Command for the request to Mikrotik Router
	Payload  []byte // Payload for the request to Mikrotik Router
	Method   string // Method for the request to Mikrotik Router
}

// URL constructs the request URL based on the method and path.
func (r *APIRequest) URL() string {
	protocol := determineProtocolFromURL(r.Host)                  // Determine the protocol from the URL
	path := r.Command                                             // Set the path to the command
	return fmt.Sprintf("%s://%s/rest/%s", protocol, r.Host, path) // Return the URL
}

// createAndExecuteRequest creates a request configuration and executes the request.
func createAndExecuteRequest(ctx context.Context, request *APIRequest) (interface{}, error) {

	// Create a request configuration
	config := requestConfig{
		URL:      request.URL(),    // Set the URL
		Method:   request.Method,   // Set the method
		Username: request.Username, // Set the username
		Password: request.Password, // Set the password
		Payload:  request.Payload,  // Set the payload
	}

	// Execute the request and return the result and error
	return makeRequest(ctx, config)
}

// Auth function now uses createAndExecuteRequest with MethodGet
func Auth(ctx context.Context, config AuthConfig) (interface{}, error) {
	// Use the createAndExecuteRequest function with MethodGet
	return createAndExecuteRequest(ctx, &APIRequest{
		Host:     config.Host,       // Adjust the host as needed
		Username: config.Username,   // Adjust the username as needed
		Password: config.Password,   // Adjust the password as needed
		Command:  "system/resource", // Adjust the command as needed
		Method:   MethodGet,         // Use MethodGet
		Payload:  nil,               // No payload needed
	})
}

// Print creates a new GET request and executes it, returning the result and error.
func Print(ctx context.Context, host, username, password, command string) (interface{}, error) {

	// Create a new APIRequest
	request := &APIRequest{
		Host:     host,      // Set the host
		Username: username,  // Set the username
		Password: password,  // Set the password
		Command:  command,   // Set the command
		Method:   MethodGet, // Set the method
	}

	// Return the result and error
	return createAndExecuteRequest(ctx, request)
}

// Add creates a new PUT request and executes it, returning the result and error.
func Add(ctx context.Context, host, username, password, command string, payload []byte) (interface{}, error) {

	// Create a new APIRequest
	request := &APIRequest{
		Host:     host,      // Set the host
		Username: username,  // Set the username
		Password: password,  // Set the password
		Command:  command,   // Set the command
		Payload:  payload,   // Set the payload
		Method:   MethodPut, // Set the method
	}

	// Return the result and error
	return createAndExecuteRequest(ctx, request)
}

// Set creates a new PATCH request and executes it, returning the result and error.
func Set(ctx context.Context, host, username, password, command string, payload []byte) (interface{}, error) {

	// Create a new APIRequest
	request := &APIRequest{
		Host:     host,        // Set the host
		Username: username,    // Set the username
		Password: password,    // Set the password
		Command:  command,     // Set the command
		Payload:  payload,     // Set the payload
		Method:   MethodPatch, // Set the method
	}

	// Return the result and error
	return createAndExecuteRequest(ctx, request)
}

// Remove creates a new DELETE request and executes it, returning the result and error.
func Remove(ctx context.Context, host, username, password, command string) (interface{}, error) {

	// Create a new APIRequest
	request := &APIRequest{
		Host:     host,         // Set the host
		Username: username,     // Set the username
		Password: password,     // Set the password
		Command:  command,      // Set the command
		Method:   MethodDelete, // Set the method
	}

	// Return the result and error
	return createAndExecuteRequest(ctx, request)
}

// Run creates a new POST request and executes it, returning the result and error.
func Run(ctx context.Context, host, username, password, command string, payload []byte) (interface{}, error) {

	// Create a new APIRequest
	request := &APIRequest{
		Host:     host,       // Set the host
		Username: username,   // Set the username
		Password: password,   // Set the password
		Command:  command,    // Set the command
		Payload:  payload,    // Set the payload
		Method:   MethodPost, // Set the method
	}

	// Return the result and error
	return createAndExecuteRequest(ctx, request)
}
