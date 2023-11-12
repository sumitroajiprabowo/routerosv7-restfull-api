// api_methods.go

package routerosv7_restfull_api

import (
	"context"
	"fmt"
)

// HTTP methods constants.
const (
	MethodGet    = "GET"
	MethodPost   = "POST"
	MethodPut    = "PUT"
	MethodPatch  = "PATCH"
	MethodDelete = "DELETE"
)

// AuthConfig is the configuration for the AuthDevice function
type AuthConfig struct {
	Host     string // Host for the request to Mikrotik Router
	Username string // Username for the request to Mikrotik Router
	Password string // Password for the request to Mikrotik Router
}

// APIRequest represents a request to the API.
type APIRequest struct {
	Host     string
	Username string
	Password string
	Command  string
	Payload  []byte
	Method   string
}

// URL constructs the request URL based on the method and path.
func (r *APIRequest) URL() string {
	protocol := determineProtocolFromURL(r.Host)
	path := r.Command
	return fmt.Sprintf("%s://%s/rest/%s", protocol, r.Host, path)
}

// createAndExecuteRequest creates a request configuration and executes the request.
func createAndExecuteRequest(ctx context.Context, request *APIRequest) (interface{}, error) {
	config := requestConfig{
		URL:      request.URL(),
		Method:   request.Method,
		Username: request.Username,
		Password: request.Password,
		Payload:  request.Payload,
	}

	// Execute the request and return the result and error
	return makeRequest(ctx, config)
}

func Auth(ctx context.Context, config AuthConfig) (interface{}, error) {
	// Determine the protocol from the URL (HTTP or HTTPS)
	protocol := determineProtocol(config.Host)

	// Create the URL for the request to Mikrotik Router
	url := fmt.Sprintf("%s://%s/rest/system/resource", protocol, config.Host)

	// Create the request configuration for the request to Mikrotik Router
	requestConfig := requestConfig{
		URL:      url,             // URL for the request to Mikrotik Router
		Method:   "GET",           // Method for the request to Mikrotik Router
		Payload:  nil,             // Payload for the request to Mikrotik Router
		Username: config.Username, // Username for the request to Mikrotik Router
		Password: config.Password, // Password for the request to Mikrotik Router
	}

	// Make the request to Mikrotik Router
	return makeRequest(ctx, requestConfig)
}

// Print creates a new GET request and executes it, returning the result and error.
func Print(ctx context.Context, host, username, password, command string) (interface{}, error) {
	request := &APIRequest{
		Host:     host,
		Username: username,
		Password: password,
		Command:  command,
		Method:   MethodGet,
	}

	return createAndExecuteRequest(ctx, request)
}

// Add creates a new PUT request and executes it, returning the result and error.
func Add(ctx context.Context, host, username, password, command string, payload []byte) (interface{}, error) {
	request := &APIRequest{
		Host:     host,
		Username: username,
		Password: password,
		Command:  command,
		Payload:  payload,
		Method:   MethodPut,
	}

	return createAndExecuteRequest(ctx, request)
}

// Set creates a new PATCH request and executes it, returning the result and error.
func Set(ctx context.Context, host, username, password, command string, payload []byte) (interface{}, error) {
	request := &APIRequest{
		Host:     host,
		Username: username,
		Password: password,
		Command:  command,
		Payload:  payload,
		Method:   MethodPatch,
	}

	return createAndExecuteRequest(ctx, request)
}

// Remove creates a new DELETE request and executes it, returning the result and error.
func Remove(ctx context.Context, host, username, password, command string) (interface{}, error) {
	request := &APIRequest{
		Host:     host,
		Username: username,
		Password: password,
		Command:  command,
		Method:   MethodDelete,
	}

	return createAndExecuteRequest(ctx, request)
}

// Run creates a new POST request and executes it, returning the result and error.
func Run(ctx context.Context, host, username, password, command string, payload []byte) (interface{}, error) {
	request := &APIRequest{
		Host:     host,
		Username: username,
		Password: password,
		Command:  command,
		Payload:  payload,
		Method:   MethodPost,
	}

	return createAndExecuteRequest(ctx, request)
}
