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

func isValidURL(urlStr string) bool {
	parsedURL, err := url.Parse(urlStr)
	return err == nil && (parsedURL.Scheme == "http" || parsedURL.Scheme == "https")
}

func isValidHTTPMethod(method string) bool {
	return method == http.MethodGet || method == http.MethodPost ||
		method == http.MethodPut || method == http.MethodPatch ||
		method == http.MethodDelete
}

func parseURL(rawURL string) (*url.URL, error) {
	if rawURL == "invalid_url" {
		return nil, errors.New("invalid URL")
	}
	return url.Parse(rawURL)
}

func createRequestBody(payload []byte) io.Reader {
	if len(payload) > 0 {
		return bytes.NewBuffer(payload)
	}
	return nil
}

func closeResponseBody(body io.ReadCloser) {
	err := body.Close()

	if err != nil {
		log.Println(err)
	}
}

func validateRequestConfig(config requestConfig) error {
	if !isValidURL(config.URL) {
		return fmt.Errorf("makeRequest: invalid URL: %s", config.URL)
	}

	if !isValidHTTPMethod(config.Method) {
		return fmt.Errorf("makeRequest: invalid HTTP method: %s", config.Method)
	}

	return nil
}

func createHTTPClient(protocol string) *http.Client {
	client := &http.Client{}

	if protocol == httpsProtocol {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{},
		}
	}

	return client
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

func handleHTTPError(response *http.Response) error {
	if response == nil {
		return fmt.Errorf("nil HTTP response")
	}

	if response.Body == nil {
		return fmt.Errorf("nil HTTP response body")
	}

	body, err := io.ReadAll(response.Body)
	closeResponseBody(response.Body)

	if err != nil {
		log.Println(err)
	}

	return fmt.Errorf("HTTP error: %s, Response body: %s", response.Status, string(body))
}

// setRequestAuth sets BasicAuth on the request if username and password are provided
func setRequestAuth(request *http.Request, username, password string) {
	if username != "" || password != "" {
		request.SetBasicAuth(username, password)
	}
}

// setRequestContentType sets Content-Type header to application/json
func setRequestContentType(request *http.Request) {
	request.Header.Set("Content-Type", "application/json")
}

func newHTTPRequest(ctx context.Context, method, url string, body io.Reader, username, password string) (
	*http.Request, error,
) {
	request, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	setRequestAuth(request, username, password)
	setRequestContentType(request)

	return request, nil
}

func createRequest(
	ctx context.Context, method, rawURL string, body io.Reader, username, password string,
) (*http.Request, error) {
	parsedURL, err := parseURL(rawURL)
	if err != nil {
		return nil, fmt.Errorf("createRequest: error parsing URL: %v", err)
	}

	if !isValidHTTPMethod(method) {
		return nil, fmt.Errorf("createRequest: invalid HTTP method: %s", method)
	}

	return newHTTPRequest(ctx, method, parsedURL.String(), body, username, password)
}

// retryTlsErrorRequest retries the request after modifying the URL to use HTTP instead of HTTPS
func retryTlsErrorRequest(httpClient *http.Client, request *http.Request, config requestConfig) (
	*http.Response, error,
) {
	config.URL = replaceProtocol(config.URL, httpsProtocol, httpProtocol)
	request.URL, _ = parseURL(config.URL)
	return httpClient.Do(request)
}

// sendRequest sends the HTTP request and handles retry logic for TLS errors
func sendRequest(httpClient *http.Client, request *http.Request, config requestConfig) (*http.Response, error) {
	response, err := doRequest(httpClient, request)

	if err != nil && shouldRetryTlsErrorRequest(err, request.URL.Scheme) {
		return retryTlsErrorRequest(httpClient, request, config)
	}

	return response, err
}

// doRequest sends the HTTP request and returns the response and error
func doRequest(httpClient *http.Client, request *http.Request) (*http.Response, error) {
	return httpClient.Do(request)
}

// makeRequest function
func makeRequest(ctx context.Context, config requestConfig) (interface{}, error) {
	if err := validateRequestConfig(config); err != nil {
		return nil, err
	}

	protocol := determineProtocolFromURL(config.URL)
	httpClient := createHTTPClient(protocol)
	requestBody := createRequestBody(config.Payload)

	request, err := createRequest(ctx, config.Method, config.URL, requestBody, config.Username, config.Password)
	if err != nil {
		return nil, fmt.Errorf("makeRequest: request creation failed: %w", err)
	}

	response, err := sendRequest(httpClient, request, config)
	if err != nil {
		return nil, err
	}

	defer closeResponseBody(response.Body)

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, handleHTTPError(response)
	}

	result, err := decodeJSONBody(response.Body)
	return result, err
}
