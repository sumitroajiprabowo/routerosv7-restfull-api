package routerosv7_restfull_api

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
)

func Delete(ctx context.Context, host, username, password, command string) error {

	// Determine the protocol from the URL (HTTP or HTTPS)
	protocol := determineProtocol(host)

	// Create the URL for the request
	url := fmt.Sprintf("%s://%s/rest/%s", protocol, host, command)

	// Create a new HTTP request based on the provided URL and method
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	// Set basic auth for the request
	req.SetBasicAuth(username, password)

	// Set the content type to JSON
	req.Header.Set("Content-Type", "application/json")

	// Client for HTTP requests
	client := &http.Client{}

	// Set the TLS configuration for the client based on the protocol (HTTP or HTTPS)
	if protocol == "https" {
		// Use a client with standard TLS verification for HTTPS
		client = &http.Client{}
	} else {
		// Use a client with insecure TLS configuration for HTTP
		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
	}

	req = req.WithContext(ctx)

	// Send the request
	_, err = client.Do(req)

	if err != nil {
		return err
	}

	return nil
}
