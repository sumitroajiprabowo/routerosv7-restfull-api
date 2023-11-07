package routerosv7_restfull_api

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

func makeRequest(ctx context.Context, url, username, password, method string, payload []byte) (interface{}, error) {
	// Determine the protocol from the URL (HTTP or HTTPS)
	protocol := "http"
	if strings.HasPrefix(url, "https://") {
		protocol = "https"
	}

	// Create a new HTTP client with a timeout using the context
	client := &http.Client{}

	// Use a client with standard TLS verification for HTTPS
	if protocol == "https" {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{},
		}
	}

	// Create a new buffer based on the provided payload
	var requestBody io.Reader

	// Check if payload is not empty
	if len(payload) > 0 {
		requestBody = bytes.NewBuffer(payload)
	}

	// Create a new HTTP request based on the provided URL and method
	request, err := http.NewRequestWithContext(ctx, method, url, requestBody)
	if err != nil {
		return nil, err
	}

	// Set basic auth for the request
	request.SetBasicAuth(username, password)

	// Send the request
	response, err := client.Do(request)

	if err != nil {
		// Handle TLS handshake error and retry with HTTP
		if strings.Contains(err.Error(), "tls: handshake failure") && protocol == "https" {
			protocol = "http"
			url = strings.Replace(url, "https://", "http://", 1)
			request.URL, _ = request.URL.Parse(url)
			response, err = client.Do(request)
		}

		if err != nil {
			return nil, err
		}
	}

	// Close response body and check for context cancellation
	select {
	case <-ctx.Done():
		// Request was canceled
		return nil, ctx.Err()
	default:
		// Continue processing
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Println(err)
			}
		}(response.Body)
	}

	// Read JSON response body
	var responseData interface{}

	// Decode JSON response body
	err = json.NewDecoder(response.Body).Decode(&responseData)
	if err != nil {
		return nil, err
	}

	return responseData, nil
}
