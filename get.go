package routerosv7_restfull_api

import (
	"context"
	"fmt"
)

func Print(ctx context.Context, host, username, password, command string) (interface{}, error) {

	// Determine the protocol from the URL (HTTP or HTTPS)
	protocol, err := determineProtocol(host)
	if err != nil {
		return nil, err
	}

	// Create the URL for the request
	url := fmt.Sprintf("%s://%s/rest/%s", protocol, host, command)

	config := requestConfig{
		URL:      url,
		Method:   "GET",
		Payload:  nil,
		Username: username,
		Password: password,
	}

	return makeRequest(ctx, config)
}
