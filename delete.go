package routerosv7_restfull_api

import (
	"context"
	"fmt"
)

// Delete is used to delete data from the router using the RESTful API.
func Delete(ctx context.Context, host, username, password, command string) (interface{}, error) {

	// Determine the protocol from the URL (HTTP or HTTPS)
	protocol := determineProtocol(host)

	// Create the URL for the request
	url := fmt.Sprintf("%s://%s/rest/%s", protocol, host, command)

	// Create the request config
	config := requestConfig{
		URL:      url,      // URL for the request
		Method:   "DELETE", // Method for the requests
		Payload:  nil,      // Payload for the request
		Username: username, // Username for the request
		Password: password, // Password for the request
	}

	// Make the request
	return makeRequest(ctx, config)
}
