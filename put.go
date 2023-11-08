package routerosv7_restfull_api

import (
	"context"
	"fmt"
)

func PutData(ctx context.Context, host, username, password, command string, payload []byte) (interface{}, error) {
	protocol := determineProtocolFromURL(host)
	url := fmt.Sprintf("%s://%s/rest/%s", protocol, host, command)
	config := requestConfig{
		URL:      url,
		Method:   "PUT",
		Payload:  payload,
		Username: username,
		Password: password,
	}

	return makeRequest(ctx, config)
}
