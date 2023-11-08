package routerosv7_restfull_api

import (
	"context"
	"fmt"
)

// AuthDeviceConfig is the configuration for the AuthDevice function
type AuthDeviceConfig struct {
	Host     string // Host for the request to Mikrotik Router
	Username string // Username for the request to Mikrotik Router
	Password string // Password for the request to Mikrotik Router
}

func Auth(ctx context.Context, config AuthDeviceConfig) (interface{}, error) {
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
