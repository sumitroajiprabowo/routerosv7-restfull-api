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

// makeRequest makes a request to the RouterOS API
func makeRequest(ctx context.Context, config requestConfig) (interface{}, error) {

	// Determine the protocol from the URL (HTTP or HTTPS)
	protocol := determineProtocolFromURL(config.URL)

	// Create the HTTP client
	httpClient := createHTTPClient(protocol)

	// Create the request body
	requestBody := createRequestBody(config.Payload)

	// Create the request object
	request, err := createRequest(ctx, config.Method, config.URL, requestBody, config.Username, config.Password)

	// Return the error if request creation failed
	if err != nil {
		return nil, err
	}

	// Make the request with the HTTP client and return the response
	response, err := httpClient.Do(request)

	// Retry the request if it failed and the protocol is HTTPS
	if err != nil {
		if shouldRetryRequest(err, protocol) {

			// Retry the request with HTTP
			config.URL = replaceProtocol(config.URL, httpsProtocol, httpProtocol)

			// Parse the URL
			request.URL, _ = request.URL.Parse(config.URL)

			// Make the request with HTTP
			response, err = httpClient.Do(request)
		}

		// Return the error if the request still failed
		if err != nil {
			return nil, err
		}
	}

	// Close the response body
	defer closeResponseBody(response.Body)

	// Handle non-2xx status codes as errors
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, handleHTTPError(response)
	}

	// Decode the JSON body
	return decodeJSONBody(response.Body)
}

// handleHTTPError handles non-2xx status codes as errors
func handleHTTPError(response *http.Response) error {

	// Read the response body
	body, err := io.ReadAll(response.Body)

	// Return the error if there is one
	if err != nil {
		return err
	}

	// Return the HTTP error
	return fmt.Errorf("HTTP error: %s, Response body: %s", response.Status, string(body))
}

// Determine the protocol from the URL (HTTP or HTTPS)
func createHTTPClient(protocol string) *http.Client {

	// Create the HTTP client
	client := &http.Client{}

	// Check if the protocol is HTTPS and set the TLS config if it's HTTPS
	if protocol == httpsProtocol {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{},
		}
	}

	// Return the HTTP client
	return client
}

// Create the request body from the payload
func createRequestBody(payload []byte) io.Reader {

	// Check if the payload is not empty
	if len(payload) > 0 {
		return bytes.NewBuffer(payload)
	}
	return nil
}

func createRequest(ctx context.Context, method, url string, body io.Reader, username, password string) (
	*http.Request, error,
) {

	// Create the request object
	request, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	// Set the basic auth credentials
	request.SetBasicAuth(username, password)

	// Set the content type to JSON
	request.Header.Set("Content-Type", "application/json")

	// Return the request object and nil error
	return request, nil
}

// Close the response body and log the error if there is one
func closeResponseBody(body io.ReadCloser) {

	// Close the response body
	err := body.Close()

	// Log the error if there is one
	if err != nil {
		log.Println(err)
	}
}

// Decode the JSON body and return the data as interface{} and nil error
func decodeJSONBody(body io.ReadCloser) (interface{}, error) {

	// Decode the JSON body to interface{}
	var responseData interface{}

	// Check if there is an error while decoding the JSON body
	if err := json.NewDecoder(body).Decode(&responseData); err != nil {
		return nil, err
	}

	// Return the data as interface{} and nil error
	return responseData, nil
}
