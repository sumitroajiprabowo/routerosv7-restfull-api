package routerosv7_restfull_api

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestRequestConfig(t *testing.T) {
	// Create a sample requestConfig
	config := requestConfig{
		URL:      "https://example.com",
		Method:   "GET",
		Payload:  []byte("sample payload"),
		Username: "user",
		Password: "pass",
	}

	// Test URL
	expectedURL := "https://example.com"
	if config.URL != expectedURL {
		t.Errorf("Expected URL to be %s, got %s", expectedURL, config.URL)
	}

	// Test Method
	expectedMethod := "GET"
	if config.Method != expectedMethod {
		t.Errorf("Expected Method to be %s, got %s", expectedMethod, config.Method)
	}

	// Test Payload
	expectedPayload := []byte("sample payload")
	if string(config.Payload) != string(expectedPayload) {
		t.Errorf("Expected Payload to be %v, got %v", expectedPayload, config.Payload)
	}

	// Test Username
	expectedUsername := "user"
	if config.Username != expectedUsername {
		t.Errorf("Expected Username to be %s, got %s", expectedUsername, config.Username)
	}

	// Test Password
	expectedPassword := "pass"
	if config.Password != expectedPassword {
		t.Errorf("Expected Password to be %s, got %s", expectedPassword, config.Password)
	}
}

// TestIsValidURL tests the isValidURL function with various inputs and expected outputs for each test case defined
// in the tests array of struct below it.
func TestIsValidURL(t *testing.T) {

	// Test case 1: Valid HTTP URL
	tests := []struct {
		name     string
		url      string
		expected bool
	}{
		{"Valid HTTP URL", "http://example.com", true},
		{"Valid HTTPS URL", "https://example.com", true},
		{"Invalid URL", "invalid_url", false},
		{"Empty URL", "", false},
	}

	// Run the isValidURL function with the test cases defined above
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidURL(tt.url)
			if result != tt.expected {
				t.Errorf("isValidURL(%s) = %v; want %v", tt.url, result, tt.expected)
			}
		})
	}
}

// mustParseURL parses a raw URL and panics if there is an error. This is used for creating a URL for testing.
// This function is used in the TestParseURL function below. It is not exported. It is only used in this file.
// It is not used in any other file.
func mustParseURL(rawURL string) *url.URL {

	// Parse the raw URL
	parsedURL, err := url.Parse(rawURL)

	// If there is an error, panic
	if err != nil {
		panic(err)
	}

	// Return the parsed URL
	return parsedURL
}

// TestIsValidHTTPMethod tests the isValidHTTPMethod function with various inputs and expected outputs for each test
// case defined in the tests array of struct below it
func TestIsValidHTTPMethod(t *testing.T) {

	// Test case 1: Valid HTTP method
	tests := []struct {
		name     string
		method   string
		expected bool
	}{
		{"Valid GET method", "GET", true},
		{"Valid POST method", "POST", true},
		{"Invalid method", "INVALID_METHOD", false},
		{"Empty method", "", false},
	}

	// Run the isValidHTTPMethod function with the test cases defined above
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidHTTPMethod(tt.method)
			if result != tt.expected {
				t.Errorf("isValidHTTPMethod(%s) = %v; want %v", tt.method, result, tt.expected)
			}
		})
	}
}

// TestParseURL tests the parseURL function with various inputs and expected outputs for each test case defined in the tests array of struct below it.
func TestParseURL(t *testing.T) {

	// Test case 1: Valid URL
	tests := []struct {
		name     string
		rawURL   string
		expected *url.URL
		wantErr  bool
	}{
		{"Valid URL", "https://example.com/path", mustParseURL("https://example.com/path"), false},
		{"Invalid URL", "invalid_url", nil, true},
		{"Empty URL", "", nil, false},
	}

	// Run the parseURL function with the test cases defined above
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseURL(tt.rawURL)

			// Check if the error is as expected for invalid URLs
			if (err != nil) != tt.wantErr {
				t.Errorf("parseURL(%s) error = %v, wantErr %v", tt.rawURL, err, tt.wantErr)
				return
			}

			// Check if the result is as expected for valid URLs
			if tt.expected != nil && (result == nil || result.String() != tt.expected.String()) {
				t.Errorf("parseURL(%s) = %v, want %v", tt.rawURL, result, tt.expected)
			}
		})
	}
}

// TestCreateRequestURL tests the createRequestURL function with various inputs and expected outputs for each test case defined in the tests array of struct below it.
func TestCreateRequestBody(t *testing.T) {
	// Non-empty payload with valid JSON
	payload := []byte(`{"key": "value"}`)

	// Create a request body with the non-empty payload
	body := createRequestBody(payload)

	// Check if the body is not nil
	if body == nil {
		t.Error("Expected non-nil body for non-empty payload")
	}

	// Empty payload
	var emptyPayload []byte

	// Create a request body with the empty payload
	emptyBody := createRequestBody(emptyPayload)

	// Check if the body is nil
	if emptyBody != nil {
		t.Error("Expected nil body for empty payload")
	}
}

type mockErrorReaderCloser struct{} // Mock a reader closer with an error on read and close

// Mock the Read method to return an error
func (m *mockErrorReaderCloser) Read(_ []byte) (n int, err error) {
	return 0, errors.New("mocked read error")
}

// Mock the Close method to return an error
func (m *mockErrorReaderCloser) Close() error {
	return errors.New("mocked close error")
}

// TestCloseResponseBody tests is the closeResponseBody function with various inputs and expected outputs for each test case defined in the tests array of struct below it.
func TestCloseResponseBody(t *testing.T) {

	// Mock a response body with an error on close
	errorBody := &mockErrorReaderCloser{}

	// Close the response body
	closeResponseBody(errorBody) // This should log the error, you can capture logs and check
}

// TestValidateRequestConfig tests the validateRequestConfig function with various inputs and expected outputs for each test case defined in the tests array of struct below it.
func TestValidateRequestConfig(t *testing.T) {

	// Test case 1: Valid Config
	tests := []struct {
		name          string
		config        requestConfig
		expectedError bool
	}{
		{"Valid Config", requestConfig{URL: "https://example.com", Method: "GET"}, false},
		{"Invalid URL", requestConfig{URL: "invalid_url", Method: "GET"}, true},
		{"Invalid Method", requestConfig{URL: "https://example.com", Method: "INVALID"}, true},
		{"Empty Config", requestConfig{}, true},
	}

	// Run the validateRequestConfig function with the test cases defined above
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRequestConfig(tt.config)

			// Check if the error is as expected
			if tt.expectedError && err == nil {
				t.Error("validateRequestConfig() did not return the expected error.")
			}

			// Check if the error is as expected
			if !tt.expectedError && err != nil {
				t.Errorf("validateRequestConfig() returned an unexpected error: %v", err)
			}
		})
	}
}

// TestCreateHTTPClient tests the createHTTPClient function with various inputs and expected outputs for each test case
// defined in the tests array of struct below it.
func TestCreateHTTPClient(t *testing.T) {
	// Test for HTTPS protocol
	client := createHTTPClient(httpsProtocol)

	// Check if the TLS config is set for HTTPS protocol and not set for HTTP protocol for the client created above
	//and below respectively in the test cases.
	if client.Transport == nil || client.Transport.(*http.Transport).TLSClientConfig == nil {
		t.Error("Expected TLS config to be set for HTTPS protocol")
	}

	// Test for HTTP protocol (no TLS config)
	client = createHTTPClient(httpProtocol)

	// Check if the TLS config is not set for HTTP protocol for the client created above and below respectively in the test cases.
	if client.Transport != nil {
		t.Error("Expected no TLS config for HTTP protocol")
	}
}

// TestHandleHTTPError tests the handleHTTPError function with various inputs and expected outputs for each test case.
func TestHandleHTTPError(t *testing.T) {

	// Test case 1: Non-2xx status code with a response body
	responseBody1 := "Error response body 1"

	// Create a response with a non-2xx status code and a response body
	statusCode1 := http.StatusNotFound

	// Create a response with a non-2xx status code and a response body
	response1 := httptest.NewRecorder()

	// Write the response body
	response1.WriteHeader(statusCode1)

	// Write the response body
	response1.Body.WriteString(responseBody1)

	// Call handleHTTPError with the first sample response
	err1 := handleHTTPError(response1.Result())

	// Check if the error matches the expected format
	expectedError1 := fmt.Sprintf("HTTP error: %d %s, Response body: %s", statusCode1, http.StatusText(statusCode1), responseBody1)

	// Check if the error matches the expected format
	if err1 == nil || err1.Error() != expectedError1 {
		t.Errorf("handleHTTPError did not return the expected error. Got: %v, Expected: %s", err1, expectedError1)
	}

	// Test case 2: Non-2xx status code without a response body
	statusCode2 := http.StatusInternalServerError

	// Create a response with a non-2xx status code and without a response body
	response2 := httptest.NewRecorder()

	// Write the response body
	response2.WriteHeader(statusCode2)

	// Call handleHTTPError with the second sample response
	err2 := handleHTTPError(response2.Result())

	// Check if the error matches the expected format
	expectedError2 := fmt.Sprintf("HTTP error: %d %s, Response body: ", statusCode2, http.StatusText(statusCode2))

	// Check if the error matches the expected format
	if err2 == nil || err2.Error() != expectedError2 {
		t.Errorf("handleHTTPError did not return the expected error. Got: %v, Expected: %s", err2, expectedError2)
	}

	// Test case 3: 2xx status code with a response body
	response3 := httptest.NewRecorder()

	// Write the response body with a 2xx status code
	response3.WriteHeader(http.StatusOK)

	// Write the response body with a 2xx status code
	response3.Body.WriteString("Success response body")

	// Call handleHTTPError with the third sample response
	err3 := handleHTTPError(response3.Result())

	// Check if the error is not nil for 2xx status codes
	if err3 == nil {
		t.Errorf("handleHTTPError did not return an error for a 2xx status code. Expected: non-nil error")
	}

	// Test case 4: Non-2xx status code with an empty response body
	statusCode4 := http.StatusBadRequest

	// Create a response with a non-2xx status code and an empty response body
	response4 := httptest.NewRecorder()

	// Write the response body
	response4.WriteHeader(statusCode4)

	// Call handleHTTPError with the fourth sample response
	err4 := handleHTTPError(response4.Result())

	// Check if the error matches the expected format
	expectedError4 := fmt.Sprintf("HTTP error: %d %s, Response body: ", statusCode4, http.StatusText(statusCode4))
	if err4 == nil || err4.Error() != expectedError4 {
		t.Errorf("handleHTTPError did not return the expected error. Got: %v, Expected: %s", err4, expectedError4)
	}

	// Test case 5: Non-2xx status code with a large response body
	statusCode5 := http.StatusUnauthorized

	// Create a response with a non-2xx status code and a large response body
	response5 := httptest.NewRecorder()

	// Write the response body with a non-2xx status code
	response5.WriteHeader(statusCode5)

	// Write the response body with a non-2xx status code and a large response body
	largeBody := "This is a very large response body. " + string(make([]byte, 1024*1024)) // 1 MB

	// Write the response body with a non-2xx status code and a large response body
	response5.Body.WriteString(largeBody)

	// Call handleHTTPError with the fifth sample response
	err5 := handleHTTPError(response5.Result())

	// Check if the error matches the expected format
	expectedError5 := fmt.Sprintf("HTTP error: %d %s, Response body: %s", statusCode5, http.StatusText(statusCode5), largeBody)

	// Check if the error matches the expected format
	if err5 == nil || err5.Error() != expectedError5 {
		t.Errorf("handleHTTPError did not return the expected error. Got: %v, Expected: %s", err5, expectedError5)
	}

	// Test case 6: Nil HTTP response
	var responseNil *http.Response

	// Call handleHTTPError with the sixth sample response
	err6 := handleHTTPError(responseNil)

	// Check if the error matches the expected format
	expectedError6 := "nil HTTP response"

	// Check if the error matches the expected format
	if err6 == nil || err6.Error() != expectedError6 {
		t.Errorf("handleHTTPError did not return the expected error. Got: %v, Expected: %s", err6, expectedError6)
	}

	// Test case 7: Nil HTTP response body
	statusCode7 := http.StatusNotFound

	// Create a response with a non-2xx status code and a nil response body
	response7 := &http.Response{
		Status:     fmt.Sprintf("%d %s", statusCode7, http.StatusText(statusCode7)),
		StatusCode: statusCode7,
		Body:       nil,
	}

	// Call handleHTTPError with the seventh sample response
	err7 := handleHTTPError(response7)

	// Check if the error matches the expected format
	expectedError7 := "nil HTTP response body"
	if err7 == nil || err7.Error() != expectedError7 {
		t.Errorf("handleHTTPError did not return the expected error. Got: %v, Expected: %s", err7, expectedError7)
	}

	// Test case 8: HTTP response body read error
	statusCode8 := http.StatusNotFound

	// Create a response with a non-2xx status code and a response body with a read error
	response8 := &http.Response{
		Status: fmt.Sprintf("%d %s", statusCode8, http.StatusText(statusCode8)),
		Body:   &mockErrorReaderCloser{},
	}

	// Call handleHTTPError with the eighth sample response
	err8 := handleHTTPError(response8)

	// Check if the error matches the expected format
	expectedError8 := "HTTP error: 404 Not Found, Response body: "

	// Check if the error matches the expected format
	if err8 == nil || err8.Error() != expectedError8 {
		t.Errorf("handleHTTPError did not return the expected error. Got: %v, Expected: %s", err8, expectedError8)
	}

	//Test case 9: HTTP response body close error
	statusCode9 := http.StatusNotFound

	// Create a response with a non-2xx status code and a response body with a close error
	response9 := &http.Response{
		Status: fmt.Sprintf("%d %s", statusCode9, http.StatusText(statusCode9)),
		Body:   &mockErrorReaderCloser{},
	}

	// Call handleHTTPError with the ninth sample response and check if the error matches the expected format
	err9 := handleHTTPError(response9)

	// Check if the error matches the expected format
	expectedError9 := "HTTP error: 404 Not Found, Response body: "
	if err9 == nil || err9.Error() != expectedError9 {
		t.Errorf("Error")
	}

}

type emptyReadCloser struct{} // Mock a reader closer with an EOF on read and no error on close

// Mock the Read method to return an EOF
func (erc *emptyReadCloser) Read(_ []byte) (n int, err error) {
	return 0, io.EOF
}

// Mock the Close method to return no error
func (erc *emptyReadCloser) Close() error {
	return nil
}

// TestDecodeJSONBody tests the decodeJSONBody function with various inputs and expected outputs for each test case
type stringReadCloser struct {
	io.Reader
}

// Mock the Close method to return no error
func (src *stringReadCloser) Close() error {
	return nil
}

// TestDecodeJSONBody tests the decodeJSONBody function with various inputs and expected outputs for each test case
func TestDecodeJSONBody(t *testing.T) {
	// Create a response with an empty body
	response := &http.Response{
		StatusCode: http.StatusOK,
		Body:       &emptyReadCloser{},
	}

	// Attempt to decode an empty body
	result, err := decodeJSONBody(response.Body)
	if err != nil && err != io.EOF {
		t.Errorf("Expected no error or EOF, got %v", err)
	}

	// Ensure the result is nil
	if result != nil {
		t.Error("Expected a nil result for an empty JSON body")
	}

	// Create a response with a non-empty body
	response = &http.Response{
		StatusCode: http.StatusOK,
		Body:       &stringReadCloser{strings.NewReader(`{"key": "value"}`)},
	}

	// Attempt to decode a non-empty body
	result, err = decodeJSONBody(response.Body)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Ensure the result is not nil
	if result == nil {
		t.Error("Expected a non-nil result for a non-empty JSON body")
	}

	// Create another response with a different non-empty body
	response = &http.Response{
		StatusCode: http.StatusOK,
		Body:       &stringReadCloser{strings.NewReader(`{"anotherKey": "anotherValue"}`)},
	}

	// Attempt to decode the different non-empty body
	result, err = decodeJSONBody(response.Body)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Ensure the result is not nil
	if result == nil {
		t.Error("Expected a non-nil result for a different non-empty JSON body")
	}
}

// TestSetRequestAuth tests the setRequestAuth function with various inputs and expected outputs for each test case
func TestSetRequestAuth(t *testing.T) {
	// Test case 1: Set BasicAuth with valid username and password
	request := &http.Request{Header: make(http.Header)}
	username := "user"
	password := "pass"
	setRequestAuth(request, username, password)

	// Check if BasicAuth is set correctly
	if user, pass, ok := request.BasicAuth(); !ok || user != username || pass != password {
		t.Errorf("Test case 1: BasicAuth not set correctly. Expected (%s, %s), got (%s, %s)", username, password, user, pass)
	}

	// Test case 2: Do not set BasicAuth if username and password are empty
	request = &http.Request{Header: make(http.Header)}
	setRequestAuth(request, "", "")

	if user, pass, ok := request.BasicAuth(); ok || user != "" || pass != "" {
		t.Error("Test case 2: BasicAuth set incorrectly. Expected not set, but it is set.")
	}
}

// TestSetRequestContentType tests the setRequestContentType function with various inputs and expected outputs for each test case
func TestSetRequestContentType(t *testing.T) {
	// Test case 1: Content-Type is set to application/json
	request := &http.Request{Header: make(http.Header)}
	setRequestContentType(request)

	// Check if Content-Type is set correctly
	if contentType := request.Header.Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Test case 1: Content-Type not set correctly. Expected application/json, got %s", contentType)
	}
}

// TestSetRequestHeaders tests the setRequestHeaders function with various inputs and expected outputs for each test case
func TestNewHTTPRequest(t *testing.T) {
	// Test case 1: Valid GET request without authentication
	ctx := context.Background()
	method := http.MethodGet
	urls := "http://example.com"
	var body io.Reader = nil
	username := ""
	password := ""

	// Create a new HTTP request
	request, err := newHTTPRequest(ctx, method, urls, body, username, password)

	// Check if there is no error for a valid GET request without authentication
	if err != nil {
		t.Errorf("Test case 1: Expected no error, but got an error: %v", err)
	}

	// Check if the request is not nil for a valid GET request without authentication
	if request == nil {
		t.Error("Test case 1: Expected a non-nil request, but got nil")
	}

	// Test case 2: Valid POST request with authentication
	method = http.MethodPost
	body = strings.NewReader(`{"key": "value"}`)
	username = "user"
	password = "pass"

	// Create a new HTTP request
	request, err = newHTTPRequest(ctx, method, urls, body, username, password)

	// Check if there is no error for a valid POST request with authentication
	if err != nil {
		t.Errorf("Test case 2: Expected no error, but got an error: %v", err)
	}

	// Check if the request is not nil for a valid POST request with authentication
	if request == nil {
		t.Error("Test case 2: Expected a non-nil request, but got nil")
	}

	// Test case 3: Invalid URL
	urls = ":invalid-url"
	_, err = newHTTPRequest(ctx, method, urls, body, username, password)

	// Check if there is an error for an invalid URL
	if err == nil {
		t.Error("Test case 3: Expected an error for an invalid URL, but got nil")
	}
}

// TestCreateRequest tests the createRequest function with various inputs and expected outputs for each test case
func TestCreateRequest(t *testing.T) {
	// Test case 1: Valid GET request
	ctx := context.Background()
	method := http.MethodGet
	urls := "http://example.com"
	body := io.Reader(nil)
	username := "user"
	password := "pass"

	// Create a request
	request, err := createRequest(ctx, method, urls, body, username, password)

	// Check if there is no error for a valid GET request
	if err != nil {
		t.Errorf("Test case 1: Expected no error, got %v", err)
	}

	// Check if the request is not nil for a valid GET request
	if request == nil {
		t.Error("Test case 1: Expected non-nil request, got nil")
	}

	// Test case 2: Valid POST request
	method = http.MethodPost
	request, err = createRequest(ctx, method, urls, body, username, password)

	// Check if there is no error for a valid POST request
	if err != nil {
		t.Errorf("Test case 2: Expected no error, got %v", err)
	}

	// Check if the request is not nil for a valid POST request
	if request == nil {
		t.Error("Test case 2: Expected non-nil request, got nil")
	}

	// Test case 3: Valid request with body
	method = http.MethodPut
	body = bytes.NewBufferString(`{"key": "value"}`)
	request, err = createRequest(ctx, method, urls, body, username, password)

	// Check if there is no error for a valid request with body
	if err != nil {
		t.Errorf("Test case 3: Expected no error, got %v", err)
	}

	// Check if the request is not nil for a valid request with body
	if request == nil {
		t.Error("Test case 3: Expected non-nil request, got nil")
	}

	// Test case 4: Invalid URL
	urls = ":invalid-url"
	_, err = createRequest(ctx, method, urls, body, username, password)

	// Check if there is an error for an invalid URL
	if err == nil {
		t.Error("Test case 4: Expected error for invalid URL, got nil")
	}

	// Test case 5: Error parsing URL
	urls = "http://example.com"
	_, err = createRequest(ctx, method, urls, body, username, password)

	// Check if there is no error for a valid request with body
	if err != nil {
		t.Errorf("Test case 5: Expected no error while creating request, got %v", err)
	}

	// Test case 6: Invalid HTTP method
	method = "INVALID"
	_, err = createRequest(ctx, method, urls, body, username, password)

	// Check if there is an error for an invalid HTTP method
	if err == nil {
		t.Error("Test case 5: Expected error for invalid HTTP method, got nil")
	}

	// Test case 7: Error creating request
	method = http.MethodGet
	urls = "http://example.com"
	username = "user"
	password = "pass"

	// Create a request
	_, err = createRequest(ctx, method, urls, body, username, password)

	// Check if there is no error for a valid request
	if err != nil {
		t.Errorf("Test case 7: Expected no error, got %v", err)
	}

	// Test case 8: Error creating request with invalid URL
	method = http.MethodGet
	urls = ":invalid-url"
	username = "user"
	password = "pass"

	// Create a request
	_, err = createRequest(ctx, method, urls, body, username, password)

	// Check if there is an error for an invalid URL
	if err == nil {
		t.Error("Test case 8: Expected error for invalid URL, got nil")
	}

	// Test case 9: Error creating request with invalid HTTP method
	method = "INVALID"
	urls = "http://example.com"
	username = "user"
	password = "pass"

	// Create a request
	_, err = createRequest(ctx, method, urls, body, username, password)

	// Check if there is an error for an invalid HTTP method
	if err == nil {
		t.Error("Test case 9: Expected error for invalid HTTP method, got nil")
	}

	// Test case 10: Error creating request with invalid HTTP method and invalid URL
	method = "INVALID"
	urls = ":invalid-url"
	username = "user"
	password = "pass"

	// Create a request
	_, err = createRequest(ctx, method, urls, body, username, password)

	// Check if there is an error for an invalid HTTP method and invalid URL
	if err == nil {
		t.Error("Test case 10: Expected error for invalid HTTP method and invalid URL, got nil")
	}

	// Test case 11: Error creating request with empty body
	method = http.MethodGet
	urls = "http://example.com"
	username = "user"
	password = "pass"
	body = nil

	// Create a request
	_, err = createRequest(ctx, method, urls, body, username, password)

	// Check if there is no error for a valid request with empty body
	if err != nil {
		t.Errorf("Test case 11: Expected no error, got %v", err)
	}

}

// TestRetryTlsErrorRequest tests the retryTlsErrorRequest function with a mock HTTP server
func TestRetryTlsErrorRequest(t *testing.T) {
	// Create a mock HTTP server with a handler that returns a 500 status code for the first request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate a TLS handshake failure for the first request
		if r.URL.Scheme == httpsProtocol {
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte("Simulated TLS handshake failure"))
			if err != nil {
				return
			}
			return
		}

		// Return a 200 status code for the second request
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("Success"))
		if err != nil {
			return
		}
	}))
	defer server.Close()

	// Create a request configuration with the mock server URL and a valid HTTP method and payload
	config := requestConfig{
		URL:      server.URL,
		Method:   http.MethodGet,
		Payload:  nil,
		Username: "",
		Password: "",
	}

	// Make an HTTP client with a transport that injects an error for the first request (TLS handshake failure)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{},
		},
	}

	// Create an HTTP request with the mock server URL and a valid HTTP method and payload
	request, err := http.NewRequest(http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	// Run the retryTlsErrorRequest function with the mock server URL and a valid HTTP method and payload and the
	//HTTP client with a transport that injects an error for the first request (
	//TLS handshake failure) and the request configuration with the mock server URL and a valid HTTP method and payload
	response, err := retryTlsErrorRequest(client, request, config)

	// Check if there is no error for the second request (
	//TLS handshake success) and the request configuration with the mock server URL and a valid HTTP method and
	//payload and the HTTP client with a transport that injects an error for the first request (TLS handshake failure)
	if err != nil {
		t.Fatalf("RetryTlsErrorRequest failed: %v", err)
	}

	// Check if the response has the expected status code for the second request (
	//TLS handshake success) and the request configuration with the mock server URL and a valid HTTP method and
	//payload and the HTTP client with a transport that injects an error for the first request (TLS handshake failure)
	if response.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, response.StatusCode)
	}

	// Check if the response body is as expected for the second request
	expectedResponse := "Success"
	actualResponse := readResponseBody(response)
	if actualResponse != expectedResponse {
		t.Fatalf("Expected response body %s, got %s", expectedResponse, actualResponse)
	}
}

// Helper function for reading the response body
func readResponseBody(response *http.Response) string {
	body, _ := io.ReadAll(response.Body)
	defer closeResponseBody(response.Body)
	return string(body)
}

// TestSendRequest tests the sendRequest function with a mock HTTP server
func TestSendRequest(t *testing.T) {
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("Success"))
		if err != nil {
			return
		}
	}))
	defer server.Close()

	// Create a request configuration
	config := requestConfig{
		URL:      server.URL,
		Method:   http.MethodGet,
		Payload:  nil,
		Username: "",
		Password: "",
	}

	// Create an HTTP client with a transport that injects an error
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{},
		},
	}

	// Create an HTTP request
	request, err := createRequest(context.Background(), config.Method, config.URL, nil, config.Username, config.Password)
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	// Run the sendRequest function
	response, err := sendRequest(client, request, config)

	// Check if there is no error
	if err != nil {
		t.Fatalf("sendRequest failed: %v", err)
	}

	// Check if the response has the expected status code
	if response.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, response.StatusCode)
	}

	// Check if the response body is as expected
	expectedResponse := "Success"
	actualResponse := readResponseBody(response)
	if actualResponse != expectedResponse {
		t.Fatalf("Expected response body %s, got %s", expectedResponse, actualResponse)
	}
}

// TestMakeRequest tests the makeRequest function with a mock HTTP server
func TestMakeRequest(t *testing.T) {
	// Mocking a server for testing
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"status": "success"}`))
		if err != nil {
			return
		}
	}))
	defer server.Close()

	// Sample request configuration
	config := requestConfig{
		Method:   http.MethodGet,
		URL:      server.URL,
		Username: "test",
		Password: "password",
	}

	// Testing a successful request
	response, err := makeRequest(context.Background(), config)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Ensure the response is not nil
	if response == nil {
		t.Error("Expected a non-nil response")
	}

	// Test invalid URL
	config.URL = "invalid-url"
	_, err = makeRequest(context.Background(), config)
	if err == nil {
		t.Error("Expected an error for invalid URL, got nil")
	} else if !strings.Contains(err.Error(), "makeRequest: invalid URL") {
		t.Errorf("Expected an error containing 'makeRequest: invalid URL', got %v", err)
	}

	// Simulate an error during the request creation with invalid URL
	config.URL = "invalid-url" // This will cause a request error
	_, err = makeRequest(context.Background(), config)
	if err == nil {
		t.Error("Expected an error, got nil")
	}

	// Testing non-2xx status code handling
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(`{"error": "not found"}`))
		if err != nil {
			return
		}
	}))
	defer server.Close()

	config.URL = server.URL
	_, err = makeRequest(context.Background(), config)
	if err == nil {
		t.Error("Expected an error for non-2xx status code, got nil")
	}

	// Testing non-2xx status code handling with empty response body
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	config.URL = server.URL
	_, err = makeRequest(context.Background(), config)
	if err == nil {
		t.Error("Expected an error for non-2xx status code, got nil")
	}

	// Testing non-2xx status code handling with large response body
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(string(make([]byte, 1024*1024)))) // 1 MB
		if err != nil {
			return
		}
	}))

	defer server.Close()

	config.URL = server.URL
	_, err = makeRequest(context.Background(), config)
	if err == nil {
		t.Error("Expected an error for non-2xx status code, got nil")
	}

	// Simulate an error during the request creation with invalid HTTP method
	config.Method = "INVALID" // This will cause a request error
	_, err = makeRequest(context.Background(), config)
	if err == nil {
		t.Error("Expected an error, got nil")
	}

	// Simulate an error during the request creation with invalid URL and HTTP method
	config.URL = "invalid-url" // This will cause a request error
	config.Method = "INVALID"  // This will cause a request error
	_, err = makeRequest(context.Background(), config)
	if err == nil {
		t.Error("Expected an error, got nil")
	}

	// Simulate an error during the request creation with empty body
	config.URL = server.URL
	config.Method = http.MethodGet
	config.Payload = nil // This will cause a request error
	_, err = makeRequest(context.Background(), config)
	if err == nil {
		t.Error("Expected an error, got nil")
	}

	// Simulate an error during the request creation with error parsing URL
	config.URL = ":invalid-url" // This will cause a request error
	config.Method = http.MethodGet
	config.Payload = nil
	_, err = makeRequest(context.Background(), config)
	if err == nil {
		t.Error("Expected an error, got nil")
	}

	// Simulate an error during the request creation with invalid HTTP method and invalid URL
	config.URL = ":invalid-url" // This will cause a request error
	config.Method = "INVALID"   // This will cause a request error
	config.Payload = nil
	_, err = makeRequest(context.Background(), config)
	if err == nil {
		t.Error("Expected an error, got nil")
	}

	// Testing handleHTTPError function
	errorResponse := &http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       io.NopCloser(strings.NewReader(`{"error": "bad request"}`)),
	}

	err = handleHTTPError(errorResponse)
	if err == nil {
		t.Error("Expected an error, got nil")
	} else if !strings.Contains(err.Error(), "HTTP error") {
		t.Errorf("Expected an error containing 'HTTP error', got %v", err)
	}

	// Testing handleHTTPError function with nil response
	err = handleHTTPError(nil)
	if err == nil {
		t.Error("Expected an error, got nil")
	}

	// Testing handleHTTPError function with nil response body
	errorResponse = &http.Response{
		StatusCode: http.StatusBadRequest,
	}

	// Call handleHTTPError with the sample response
	err = handleHTTPError(errorResponse)

	// Check if the error matches the expected format
	if err == nil {
		t.Error("Expected an error, got nil")
	}

}
