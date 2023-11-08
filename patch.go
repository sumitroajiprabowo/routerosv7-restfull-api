package routerosv7_restfull_api

import (
	"context"
	"fmt"
)

func Patch(ctx context.Context, host, username, password, command string, payload []byte) (interface{}, error) {

	// Determine the protocol from the URL (HTTP or HTTPS)
	protocol := determineProtocol(host)

	// Create the URL for the request
	url := fmt.Sprintf("%s://%s/rest/%s", protocol, host, command)

	config := requestConfig{
		URL:      url,
		Method:   "PATCH",
		Payload:  payload,
		Username: username,
		Password: password,
	}

	return makeRequest(ctx, config)
}
