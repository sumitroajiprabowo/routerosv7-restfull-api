package routerosv7_restfull_api

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func AddData(ctx context.Context, host, username, password, command string, payload []byte) (interface{}, error) {

	// Determine the protocol from the URL (HTTP or HTTPS)
	protocol := determineProtocol(host)

	// Set the URL for the request
	url := fmt.Sprintf("%s://%s/rest/%s/add", protocol, host, command)

	// Create a new HTTP request based on the provided URL and method
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
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
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// Close response body
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)

	// Read JSON response body into a map
	var response map[string]interface{}

	// Decode JSON response body
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
