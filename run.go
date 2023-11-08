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

// Requester is the interface that wraps the Do method.
type Requester interface {
	Do(ctx context.Context) (interface{}, error)
}

// Do send the API request and returns the response.
func (r *APIRequest) Do(ctx context.Context) (interface{}, error) {
	config := requestConfig{
		URL:      r.URL(),
		Method:   r.Method,
		Payload:  r.Payload,
		Username: r.Username,
		Password: r.Password,
	}
	return makeRequest(ctx, config)
}

// Print creates a new GET request.
func Print(host, username, password, command string) Requester {
	return &APIRequest{
		Host:     host,
		Username: username,
		Password: password,
		Command:  command,
		Method:   MethodGet,
	}
}

// Add creates a new PUT request.
func Add(host, username, password, command string, payload []byte) Requester {
	return &APIRequest{
		Host:     host,
		Username: username,
		Password: password,
		Command:  command,
		Payload:  payload,
		Method:   MethodPut,
	}
}

// Set creates a new PATCH request.
func Set(host, username, password, command string, payload []byte) Requester {
	return &APIRequest{
		Host:     host,
		Username: username,
		Password: password,
		Command:  command,
		Payload:  payload,
		Method:   MethodPatch,
	}
}

// Delete creates a new DELETE request.
func Delete(host, username, password, command string) Requester {
	return &APIRequest{
		Host:     host,
		Username: username,
		Password: password,
		Command:  command,
		Method:   MethodDelete,
	}
}

// Command creates a new POST request. All the API features are available through the POST method
func Command(host, username, password, command string, payload []byte) Requester {
	return &APIRequest{
		Host:     host,
		Username: username,
		Password: password,
		Command:  command,
		Payload:  payload,
		Method:   MethodPost,
	}
}
