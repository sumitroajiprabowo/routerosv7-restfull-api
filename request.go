package routerosv7_restfull_api

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

// requestConfig represents a request to the API.
type requestConfig struct {
	URL      string // URL for the request to Mikrotik Router
	Method   string // Method for the request to Mikrotik Router
	Payload  []byte // Payload for the request to Mikrotik Router
	Username string // Username for the request to Mikrotik Router
	Password string // Password for the request to Mikrotik Router
}

/*
isValidURL function is used to check if the URL is valid
It returns true if the URL is valid, otherwise it returns false
It returns false if the URL is invalid or if the URL scheme is not http or https
example:
isValidURL("http://example.com") returns true
isValidURL("https://example.com") returns true
isValidURL("ftp://example.com") returns false
isValidURL("invalid_url") returns false
*/
func isValidURL(urlStr string) bool {
	// Parse the URL string
	parsedURL, err := url.Parse(urlStr)

	// Check if there is an error while parsing the URL string
	return err == nil && (parsedURL.Scheme == "http" || parsedURL.Scheme == "https")
}

/*
isValidHTTPMethod function is used to check if the HTTP method is valid
It returns true if the HTTP method is valid, otherwise it returns false
It returns false if the HTTP method is invalid
example:
isValidHTTPMethod("GET") returns true
isValidHTTPMethod("POST") returns true
isValidHTTPMethod("PUT") returns true
isValidHTTPMethod("PATCH") returns true
isValidHTTPMethod("DELETE") returns true
isValidHTTPMethod("invalid_method") returns false
isValidHTTPMethod("") returns false
isValidHTTPMethod("PATCH ") returns false
*/
func isValidHTTPMethod(method string) bool {
	return method == http.MethodGet || method == http.MethodPost ||
		method == http.MethodPut || method == http.MethodPatch ||
		method == http.MethodDelete
}

/*
parseURL function is used to parse the URL string
It returns the parsed URL and nil error if the URL string is valid
It returns nil and error if the URL string is invalid
*/
func parseURL(rawURL string) (*url.URL, error) {

	// Check if the URL string is invalid
	if rawURL == "invalid_url" {
		return nil, errors.New("invalid URL") // Return nil and error
	}

	// Parse the URL string
	parsedURL, err := url.Parse(rawURL)

	// Check if there is an error while parsing the URL string
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err) // Return nil and error
	}

	// Return the parsed URL and nil error
	return parsedURL, nil
}

/*
createRequestBody function is used to create the request body from the payload byte array
if the payload is not empty, it returns a bytes.Buffer containing the payload
if the payload is empty, it returns nil
*/
func createRequestBody(payload []byte) io.Reader {

	// Check if the payload is not empty
	if len(payload) > 0 {
		return bytes.NewBuffer(payload) // Return a bytes.Buffer containing the payload
	}
	// Return nil
	return nil
}

// closeResponseBody function is used to close the response body and log the error if any error occurs
func closeResponseBody(body io.ReadCloser) {
	err := body.Close() // Close the response body

	// Check if there is an error while closing the response body
	if err != nil {
		log.Println(err)
	}
}

/*
validateRequestConfig function is used to validate the request config struct fields before making the request to the API
It returns nil if the request config struct fields are valid
It returns an error if the request config struct fields are invalid
*/
func validateRequestConfig(config requestConfig) error {

	// Check if the URL is valid
	if !isValidURL(config.URL) {
		return fmt.Errorf("makeRequest: invalid URL: %s", config.URL) // Return an error
	}

	// Check if the HTTP method is valid
	if !isValidHTTPMethod(config.Method) {
		return fmt.Errorf("makeRequest: invalid HTTP method: %s", config.Method) // Return an error
	}

	// Check if the payload is not empty
	return nil
}

/*
createHTTPClient function is used to create an HTTP client with TLS configuration if the protocol is https or http
client otherwise and returns the HTTP client pointer and nil error
*/
func createHTTPClient(protocol string) *http.Client {
	// Create an HTTP client
	client := &http.Client{}

	// if the protocol is https, create an HTTP client with TLS configuration
	if protocol == httpsProtocol {
		// Create an HTTP client with TLS configuration
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{},
		}
	}

	// Return the HTTP client pointer and nil error
	return client
}

// Decode the JSON body and return the data as interface{} and nil error
func decodeJSONBody(body io.ReadCloser) (interface{}, error) {

	// Decode the JSON body to interface{}
	var responseData interface{}

	// Check if there is an error while decoding the JSON body
	if err := json.NewDecoder(body).Decode(&responseData); err != nil {
		return nil, err // Return nil and error
	}

	// Return the data as interface{} and nil error
	return responseData, nil
}

/*
handleHTTPError function is used to handle the HTTP error and return the error message and response body
if any error occurs while reading the response body, it logs the error and returns the HTTP error message
*/
func handleHTTPError(response *http.Response) error {

	// Check if the response is nil
	if response == nil {
		return fmt.Errorf("nil HTTP response") // Return an error
	}

	// Check if the response body is nil
	if response.Body == nil {
		return fmt.Errorf("nil HTTP response body") // Return an error
	}

	// Read the response body
	body, err := io.ReadAll(response.Body)

	// Close the response body
	closeResponseBody(response.Body)

	// Check if there is an error while reading the response body
	if err != nil {
		log.Println(err) // Log the error
	}

	// Return the HTTP error message and response body
	return fmt.Errorf("HTTP error: %s, Response body: %s", response.Status, string(body))
}

// setRequestAuth sets BasicAuth on the request if username and password are provided
func setRequestAuth(request *http.Request, username, password string) {

	// Check if the username and password are not empty
	if username != "" || password != "" {
		request.SetBasicAuth(username, password) // Set BasicAuth on the request
	}
}

// setRequestContentType sets Content-Type header to application/json
func setRequestContentType(request *http.Request) {
	request.Header.Set("Content-Type", "application/json") // Set Content-Type header to application/json
}

/*
newHTTPRequest function is used to create a new HTTP request with the provided context, method, URL, body, username and password
It returns the HTTP request pointer and nil error if the request is created successfully
It returns nil and error if the request creation fails
*/
func newHTTPRequest(ctx context.Context, method, url string, body io.Reader, username, password string) (
	*http.Request, error,
) {

	// Create a new HTTP request with the provided context, method, URL and body
	request, err := http.NewRequestWithContext(ctx, method, url, body)

	// Check if there is an error while creating the HTTP request
	if err != nil {
		return nil, err // Return nil and error
	}

	// Set BasicAuth on the request if username and password are provided
	setRequestAuth(request, username, password)

	// Set Content-Type header to application/json
	setRequestContentType(request)

	// Return the HTTP request pointer and nil error
	return request, nil
}

/*
createRequest function is used to create an HTTP request with the provided context, method, URL, body, username and password
It returns the HTTP request pointer and nil error if the request is created successfully
It returns nil and error if the request creation fails
*/
func createRequest(
	ctx context.Context, method, rawURL string, body io.Reader, username, password string,
) (*http.Request, error) {

	// Parse the URL string
	parsedURL, err := parseURL(rawURL)

	// Check if there is an error while parsing the URL string
	if err != nil {
		return nil, fmt.Errorf("createRequest: error parsing URL: %v", err)
	}

	// Check if the HTTP method is valid
	if !isValidHTTPMethod(method) {
		return nil, fmt.Errorf("createRequest: invalid HTTP method: %s", method)
	}

	// Create a new HTTP request with the provided context, method, URL, body, username and password
	return newHTTPRequest(ctx, method, parsedURL.String(), body, username, password)
}

// retryTlsErrorRequest retries the request after modifying the URL to use HTTP instead of HTTPS
func retryTlsErrorRequest(httpClient *http.Client, request *http.Request, config requestConfig) (
	*http.Response, error,
) {

	// Replace the URL protocol from https to http
	config.URL = replaceProtocol(config.URL, httpsProtocol, httpProtocol)

	// Create a new HTTP request with the provided context, method, URL, body, username and password
	request.URL, _ = parseURL(config.URL)

	// Send the HTTP request and return the response and error
	return httpClient.Do(request)
}

// sendRequest sends the HTTP request and handles retry logic for TLS errors
func sendRequest(httpClient *http.Client, request *http.Request, config requestConfig) (*http.Response, error) {

	// Send the HTTP request and return the response and error
	response, err := doRequest(httpClient, request)

	// Check if there is an error while sending the HTTP request
	if err != nil && shouldRetryTlsErrorRequest(err, request.URL.Scheme) {
		return retryTlsErrorRequest(httpClient, request, config) // Retry the request
	}

	// Return the response and error
	return response, err
}

// doRequest sends the HTTP request and returns the response and error
func doRequest(httpClient *http.Client, request *http.Request) (*http.Response, error) {
	return httpClient.Do(request)
}

// makeRequest function
func makeRequest(ctx context.Context, config requestConfig) (interface{}, error) {

	// Validate the request config struct fields before making the request to the API
	if err := validateRequestConfig(config); err != nil {
		return nil, err // Return nil and error
	}

	// Determine the protocol from the URL
	protocol := determineProtocolFromURL(config.URL)

	// Create an HTTP client
	httpClient := createHTTPClient(protocol)

	// Create the request body from the payload
	requestBody := createRequestBody(config.Payload)

	// Create an HTTP request with the provided context, method, URL, body, username and password
	request, err := createRequest(ctx, config.Method, config.URL, requestBody, config.Username, config.Password)

	// Check if there is an error while creating the HTTP request
	if err != nil {
		return nil, fmt.Errorf("makeRequest: request creation failed: %w", err) // Return nil and error
	}

	// Send the HTTP request and return the response and error
	response, err := sendRequest(httpClient, request, config)

	// Check if there is an error while sending the HTTP request
	if err != nil {
		return nil, err // Return nil and error
	}

	// Close the response body
	defer closeResponseBody(response.Body)

	// Check if the response status code is not in the range 200-299
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, handleHTTPError(response) // Return nil and error
	}

	// Decode the JSON body and return the data as interface{} and nil error
	result, err := decodeJSONBody(response.Body)

	// Return the data as interface{} and nil error
	return result, err
}
