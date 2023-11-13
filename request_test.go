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
		URL:      "https://example.com",    // Set the URL
		Method:   "GET",                    // Set the method
		Payload:  []byte("sample payload"), // Set the payload
		Username: "user",                   // Set the username
		Password: "pass",                   // Set the password
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
		name     string // Test case name
		url      string // URL
		expected bool   // Expected result
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
		name     string // Test case name
		method   string // HTTP method
		expected bool   // Expected result
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
		name     string   // Test case name
		rawURL   string   // Raw URL
		expected *url.URL // Expected parsed URL
		wantErr  bool     // Expected error
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
	return 0, errors.New("mocked read error") // Mock the Read method to return an error
}

// Mock the Close method to return an error
func (m *mockErrorReaderCloser) Close() error {
	return errors.New("mocked close error") // Mock the Close method to return an error
}

// TestCloseResponseBody tests is the closeResponseBody function with various inputs and expected outputs for each test case defined in the tests array of struct below it.
func TestCloseResponseBody(t *testing.T) {

	// Mock a response body with an error on close
	errorBody := &mockErrorReaderCloser{} // Mock a response body with an error on close

	// Close the response body
	closeResponseBody(errorBody) // This should log the error, you can capture logs and check
}

// TestValidateRequestConfig tests the validateRequestConfig function with various inputs and expected outputs for each test case defined in the tests array of struct below it.
func TestValidateRequestConfig(t *testing.T) {

	// Test case 1: Valid Config
	tests := []struct {
		name          string        // Test case name
		config        requestConfig // Request config
		expectedError bool          // Expected error
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
		t.Error("Expected no TLS config for HTTP protocol") // Check if the TLS config is not set for HTTP protocol for the client created above and below respectively in the test cases.
	}
}

// TestHandleHTTPError_Non2xxWithBody tests the handleHTTPError function with a non-2xx status code and a response body
func TestHandleHTTPError_Non2xxWithBody(t *testing.T) {
	responseBody := "Error response body"   // Set the response body to a non-empty string
	statusCode := http.StatusNotFound       // Set the status code to 404
	response := httptest.NewRecorder()      // Create a new recorder
	response.WriteHeader(statusCode)        // Set the status code to 404
	response.Body.WriteString(responseBody) // Write the response body to the response

	err := handleHTTPError(response.Result()) // Handle the HTTP error

	// Check if the error is as expected
	expectedError := fmt.Sprintf("HTTP error: %d %s, Response body: %s", statusCode, http.StatusText(statusCode), responseBody)

	// TestHandleHTTPError_Non2xxWithBody tests the handleHTTPError function with a non-2xx status code and a response body
	assertErrorEquality(t, err, expectedError)
}

// TestHandleHTTPError_Non2xxWithoutBody tests the handleHTTPError function with a non-2xx status code and no response body
func TestHandleHTTPError_Non2xxWithoutBody(t *testing.T) {
	statusCode := http.StatusInternalServerError // Set the status code to 500
	response := httptest.NewRecorder()           // Create a new recorder
	response.WriteHeader(statusCode)             // Set the status code to 500

	err := handleHTTPError(response.Result()) // Handle the HTTP error

	// Check if the error is as expected
	expectedError := fmt.Sprintf("HTTP error: %d %s, Response body: ", statusCode, http.StatusText(statusCode))

	// TestHandleHTTPError_Non2xxWithoutBody tests the handleHTTPError function with a non-2xx status code and no response body
	assertErrorEquality(t, err, expectedError)
}

// TestHandleHTTPError_Non2xxEmptyBody tests the handleHTTPError function with a non-2xx status code and an empty response body
func TestHandleHTTPError_Non2xxEmptyBody(t *testing.T) {
	statusCode := http.StatusBadRequest // Set the status code to 400
	response := httptest.NewRecorder()  // Create a new recorder
	response.WriteHeader(statusCode)    // Set the status code to 400

	err := handleHTTPError(response.Result()) // Handle the HTTP error

	// Check if the error is as expected
	expectedError := fmt.Sprintf("HTTP error: %d %s, Response body: ", statusCode, http.StatusText(statusCode))

	// TestHandleHTTPError_Non2xxEmptyBody tests the handleHTTPError function with a non-2xx status code and an empty response body
	assertErrorEquality(t, err, expectedError)
}

// TestHandleHTTPError_Non2xxLargeBody tests the handleHTTPError function with a non-2xx status code and a large response body
func TestHandleHTTPError_Non2xxLargeBody(t *testing.T) {
	statusCode := http.StatusUnauthorized                                                 // Set the status code to 401
	response := httptest.NewRecorder()                                                    // Create a new recorder
	response.WriteHeader(statusCode)                                                      // Set the status code to 401
	largeBody := "This is a very large response body. " + string(make([]byte, 1024*1024)) // 1 MB
	response.Body.WriteString(largeBody)                                                  // Write the large body to the response

	err := handleHTTPError(response.Result()) // Handle the HTTP error

	// Check if the error is as expected
	expectedError := fmt.Sprintf("HTTP error: %d %s, Response body: %s", statusCode, http.StatusText(statusCode), largeBody)

	// TestHandleHTTPError_Non2xxLargeBody tests the handleHTTPError function with a non-2xx status code and a large response body
	assertErrorEquality(t, err, expectedError)
}

// TestHandleHTTPError_NilResponseBody tests the handleHTTPError function with a nil response body
func TestHandleHTTPError_NilResponseBody(t *testing.T) {
	statusCode := http.StatusNotFound // Set the status code to 404

	// Create a nil response body
	response := &http.Response{
		Status:     fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode)), // Set the status code to 404
		StatusCode: statusCode,                                                    // Set the status code to 404
		Body:       nil,                                                           // Set the body to nil
	}

	// Handle the HTTP error
	err := handleHTTPError(response)

	// Check if the error is as expected
	expectedError := "nil HTTP response body"

	// TestHandleHTTPError_NilResponseBody tests the handleHTTPError function with a nil response body
	assertErrorEquality(t, err, expectedError)
}

func TestHandleHTTPError_ResponseBodyReadError(t *testing.T) {
	statusCode := http.StatusNotFound // Set the status code to 404

	// Create a response with a mock error reader closer
	response := &http.Response{
		Status: fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode)), // Set the status code to 404
		Body:   &mockErrorReaderCloser{},                                      // Set the body to a mock error reader closer
	}

	// Handle the HTTP error
	err := handleHTTPError(response)

	// Check if the error is as expected
	expectedError := "HTTP error: 404 Not Found, Response body: "

	// TestHandleHTTPError_ResponseBodyReadError tests the handleHTTPError function with a response body that returns an error on read
	assertErrorEquality(t, err, expectedError)
}

// TestHandleHTTPError_ResponseBodyCloseError tests the handleHTTPError function with a response body that returns an error on close
func TestHandleHTTPError_ResponseBodyCloseError(t *testing.T) {
	statusCode := http.StatusNotFound // Set the status code to 404

	// Create a response with a mock error reader closer
	response := &http.Response{
		Status: fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode)), // Set the status code to 404
		Body:   &mockErrorReaderCloser{},                                      // Set the body to a mock error reader closer
	}

	// Handle the HTTP error
	err := handleHTTPError(response)

	// Check if the error is as expected
	expectedError := "HTTP error: 404 Not Found, Response body: "

	// TestHandleHTTPError_ResponseBodyCloseError tests the handleHTTPError function with a response body that returns an error on close
	assertErrorEquality(t, err, expectedError)
}

// Helper functions assertErrorEquality is used to check if the error is as expected
func assertErrorEquality(t *testing.T, actualError error, expectedError string) {

	// Check if the error is as expected
	if actualError == nil || actualError.Error() != expectedError {
		t.Errorf("Expected error: %s, Got: %v", expectedError, actualError) // Check if the error is as expected
	}
}

type emptyReadCloser struct{} // Mock a reader closer with an EOF on read and no error on close

// Mock the Read method to return an EOF
func (erc *emptyReadCloser) Read(_ []byte) (n int, err error) {
	return 0, io.EOF // Mock the Read method to return an EOF
}

// Mock the Close method to return no error
func (erc *emptyReadCloser) Close() error {
	return nil // Mock the Close method to return no error
}

// TestDecodeJSONBody tests the decodeJSONBody function with various inputs and expected outputs for each test case
type stringReadCloser struct {
	io.Reader // Mock a reader closer with a string reader
}

// Mock the Close method to return no error
func (src *stringReadCloser) Close() error {
	return nil // Mock the Close method to return no error
}

// TestDecodeJSONBody tests the decodeJSONBody function with various inputs and expected outputs for each test case
func TestDecodeJSONBody(t *testing.T) {
	// Create a response with an empty body
	response := &http.Response{
		StatusCode: http.StatusOK,      // Set the status code to 200
		Body:       &emptyReadCloser{}, // Set the body to an empty reader closer
	}

	// Attempt to decode an empty body
	result, err := decodeJSONBody(response.Body)

	// Ensure there is no error or EOF
	if err != nil && err != io.EOF {
		t.Errorf("Expected no error or EOF, got %v", err) // Attempt to decode an empty body
	}

	// Ensure the result is nil
	if result != nil {
		t.Error("Expected a nil result for an empty JSON body") // Ensure the result is nil
	}

	// Create a response with a non-empty body
	response = &http.Response{
		StatusCode: http.StatusOK,                                            // Set the status code to 200
		Body:       &stringReadCloser{strings.NewReader(`{"key": "value"}`)}, // Set the body to a non-empty JSON string
	}

	// Attempt to decode a non-empty body
	result, err = decodeJSONBody(response.Body)

	// Ensure there is no error or EOF
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
	request := &http.Request{Header: make(http.Header)} // Create a new request
	username := "user"                                  // Set the username to a non-empty value
	password := "pass"                                  // Set the password to a non-empty value
	setRequestAuth(request, username, password)         // Set the BasicAuth

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

	// Set the Content-Type to application/json
	setRequestContentType(request)

	// Check if Content-Type is set correctly
	if contentType := request.Header.Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Test case 1: Content-Type not set correctly. Expected application/json, got %s", contentType)
	}
}

// TestSetRequestHeaders tests the setRequestHeaders function with various inputs and expected outputs for each test case
func TestNewHTTPRequest(t *testing.T) {
	// Test case 1: Valid GET request without authentication
	ctx := context.Background()  // Create a new context
	method := http.MethodGet     // Set the method to GET
	urls := "http://example.com" // Set the URL to a valid URL
	var body io.Reader = nil     // Set the body to nil
	username := ""               // Set the username to an empty string
	password := ""               // Set the password to an empty string

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
	method = http.MethodPost                     // Change the method to POST
	body = strings.NewReader(`{"key": "value"}`) // Change the body to a non-nil value
	username = "user"                            // Change the username to a non-empty value
	password = "pass"                            // Change the password to a non-empty value

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

	// Create a new HTTP request with an invalid URL
	_, err = newHTTPRequest(ctx, method, urls, body, username, password)

	// Check if there is an error for an invalid URL
	if err == nil {
		t.Error("Test case 3: Expected an error for an invalid URL, but got nil")
	}
}

// TestCreateRequestValidGet tests the createRequest function for a valid GET request
func TestCreateRequestValidGet(t *testing.T) {
	ctx := context.Background() // Create a new context
	method := http.MethodGet    // Set the method to GET
	url := "http://example.com" // Set the URL to a valid URL
	body := io.Reader(nil)      // Set the body to nil
	username := "user"          // Set the username to a non-empty value
	password := "pass"          // Set the password to a non-empty value

	// Create a new HTTP request
	request, err := createRequest(ctx, method, url, body, username, password)

	// Check if there is no error for a valid GET request
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Check if the request is not nil for a valid GET request
	if request == nil {
		t.Error("Expected non-nil request, got nil")
	}
}

// TestCreateRequestValidPost tests the createRequest function for a valid POST request
func TestCreateRequestValidPost(t *testing.T) {
	ctx := context.Background() // Create a new context
	method := http.MethodPost   // Set the method to POST
	url := "http://example.com" // Set the URL to a valid URL
	body := io.Reader(nil)      // Set the body to nil
	username := "user"          // Set the username to a non-empty value
	password := "pass"          // Set the password to a non-empty value

	// Create a new HTTP request with a valid POST request
	request, err := createRequest(ctx, method, url, body, username, password)

	// Check if there is no error for a valid POST request
	if err != nil {
		t.Errorf("Expected no error, got %v", err) // Check if there is no error for a valid POST request
	}

	// Check if the request is not nil for a valid POST request
	if request == nil {
		t.Error("Expected non-nil request, got nil") // Check if the request is not nil for a valid POST request
	}
}

// TestCreateRequestValidWithBody tests the createRequest function for a valid request with a body
func TestCreateRequestValidWithBody(t *testing.T) {
	ctx := context.Background()                       // Create a new context
	method := http.MethodPut                          // Set the method to PUT
	url := "http://example.com"                       // Set the URL to a valid URL
	body := bytes.NewBufferString(`{"key": "value"}`) // Set the body to a non-nil value
	username := "user"                                // Set the username to a non-empty value
	password := "pass"                                // Set the password to a non-empty value

	// Create a new HTTP request with a valid request with a body
	request, err := createRequest(ctx, method, url, body, username, password)

	// Check if there is no error for a valid request with a body
	if err != nil {
		t.Errorf("Expected no error, got %v", err) // Check if there is no error for a valid request with a body
	}

	// Check if the request is not nil for a valid request with a body
	if request == nil {
		t.Error("Expected non-nil request, got nil") // Check if the request is not nil for a valid request with a body
	}
}

// TestCreateRequestInvalidURL tests the createRequest function for an invalid URL
func TestCreateRequestInvalidURL(t *testing.T) {
	ctx := context.Background() // Create a new context
	method := http.MethodPut    // Set the method to PUT
	url := ":invalid-url"       // Set the URL to an invalid URL
	body := io.Reader(nil)      // Set the body to nil
	username := "user"          // Set the username to a non-empty value
	password := "pass"          // Set the password to a non-empty value

	// Create a new HTTP request with an invalid URL
	_, err := createRequest(ctx, method, url, body, username, password)

	// Check if there is an error for an invalid URL
	if err == nil {
		t.Error("Expected error for invalid URL, got nil")
	}
}

// TestCreateRequestErrorParsingURL tests the createRequest function for an error parsing the URL
func TestCreateRequestErrorParsingURL(t *testing.T) {
	ctx := context.Background() // Create a new context
	method := http.MethodGet    // Set the method to GET
	url := "http://example.com" // Set the URL to a valid URL
	body := io.Reader(nil)      // Set the body to nil
	username := "user"          // Set the username to a non-empty value
	password := "pass"          // Set the password to a non-empty value

	// Create a new HTTP request with an error parsing the URL
	_, err := createRequest(ctx, method, url, body, username, password)

	// Check if there is an error for an error parsing the URL
	if err != nil {
		t.Errorf("Expected no error while creating request, got %v", err) // Check if there is an error for an error parsing the URL
	}
}

// TestCreateRequestInvalidHTTPMethod tests the createRequest function for an invalid HTTP method
func TestCreateRequestInvalidHTTPMethod(t *testing.T) {
	ctx := context.Background() // Create a new context
	method := "INVALID"         // Set the method to an invalid HTTP method
	url := "http://example.com" // Set the URL to a valid URL
	body := io.Reader(nil)      // Set the body to nil
	username := "user"          // Set the username to a non-empty value
	password := "pass"          // Set the password to a non-empty value

	// Create a new HTTP request with an invalid HTTP method
	_, err := createRequest(ctx, method, url, body, username, password)

	// Check if there is an error for an invalid HTTP method
	if err == nil {
		t.Error("Expected error for invalid HTTP method, got nil") // Check if there is an error for an invalid HTTP method
	}
}

// TestCreateRequestErrorCreatingRequest tests the createRequest function for an error creating the request
func TestCreateRequestErrorCreatingRequest(t *testing.T) {
	ctx := context.Background() // Create a new context
	method := http.MethodGet    // Set the method to GET
	url := "http://example.com" // Set the URL to a valid URL
	body := io.Reader(nil)      // Set the body to nil
	username := "user"          // Set the username to a non-empty value
	password := "pass"          // Set the password to a non-empty value

	// Create a new HTTP request with an error creating the request
	_, err := createRequest(ctx, method, url, body, username, password)

	// Check if there is an error for an error creating the request
	if err != nil {
		t.Errorf("Expected no error, got %v", err) // Check if there is an error for an error creating the request
	}
}

// TestCreateRequestErrorCreatingRequestWithInvalidURL tests the createRequest function for an error creating the request with an invalid URL
func TestCreateRequestErrorCreatingRequestWithInvalidURL(t *testing.T) {
	ctx := context.Background() // Create a new context
	method := http.MethodGet    // Set the method to GET
	url := ":invalid-url"       // Set the URL to an invalid URL
	body := io.Reader(nil)      // Set the body to nil
	username := "user"          // Set the username to a non-empty value
	password := "pass"          // Set the password to a non-empty value

	// Create a new HTTP request with an error creating the request with an invalid URL
	_, err := createRequest(ctx, method, url, body, username, password)

	// Check if there is an error for an error creating the request with an invalid URL
	if err == nil {
		t.Error("Expected error for invalid URL, got nil") // Check if there is an error for an error creating the request with an invalid URL
	}
}

// TestCreateRequestErrorCreatingRequestWithInvalidHTTPMethod tests the createRequest function for an error creating the request with an invalid HTTP method
func TestCreateRequestErrorCreatingRequestWithInvalidHTTPMethod(t *testing.T) {
	ctx := context.Background() // Create a new context
	method := "INVALID"         // Set the method to an invalid HTTP method
	url := "http://example.com" // Set the URL to a valid URL
	body := io.Reader(nil)      // Set the body to nil
	username := "user"          // Set the username to a non-empty value
	password := "pass"          // Set the password to a non-empty value

	// Create a new HTTP request with an error creating the request with an invalid HTTP method
	_, err := createRequest(ctx, method, url, body, username, password)

	// Check if there is an error for an error creating the request with an invalid HTTP method
	if err == nil {
		t.Error("Expected error for invalid HTTP method, got nil") // Check if there is an error for an error creating the request with an invalid HTTP method
	}
}

// TestCreateRequestErrorCreatingRequestWithInvalidHTTPMethodAndURL tests the createRequest function for an error creating the request with an invalid HTTP method and invalid URL
func TestCreateRequestErrorCreatingRequestWithInvalidHTTPMethodAndURL(t *testing.T) {
	ctx := context.Background() // Create a new context
	method := "INVALID"         // Set the method to an invalid HTTP method
	url := ":invalid-url"       // Set the URL to an invalid URL
	body := io.Reader(nil)      // Set the body to nil
	username := "user"          // Set the username to a non-empty value
	password := "pass"          // Set the password to a non-empty value

	// Create a new HTTP request with an error creating the request with an invalid HTTP method and invalid URL
	_, err := createRequest(ctx, method, url, body, username, password)

	// Check if there is an error for an error creating the request with an invalid HTTP method and invalid URL
	if err == nil {
		t.Error("Expected error for invalid HTTP method and invalid URL, got nil") // Check if there is an error for an error creating the request with an invalid HTTP method and invalid URL
	}
}

// TestCreateRequestNoErrorWithEmptyBody tests the createRequest function for no error with an empty body
func TestCreateRequestNoErrorWithEmptyBody(t *testing.T) {
	ctx := context.Background() // Create a new context
	method := http.MethodGet    // Set the method to GET
	url := "http://example.com" // Set the URL to a valid URL
	body := io.Reader(nil)      // Set the body to nil
	username := "user"          // Set the username to a non-empty value
	password := "pass"          // Set the password to a non-empty value

	// Create a new HTTP request with no error with an empty body
	_, err := createRequest(ctx, method, url, body, username, password)

	// Check if there is no error for no error with an empty body
	if err != nil {
		t.Errorf("Expected no error, got %v", err) // Check if there is no error for no error with an empty body
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

		// Return a success message for the second request
		_, err := w.Write([]byte("Success")) // Write the response body
		if err != nil {
			return
		}
	}))
	defer server.Close()

	// Create a request configuration with the mock server URL and a valid HTTP method and payload
	config := requestConfig{
		URL:      server.URL,     // URL of the mock server
		Method:   http.MethodGet, // Default to GET method
		Payload:  nil,            // No payload
		Username: "",             // No username
		Password: "",             // No password
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

	// Check if the response body is as expected for the second request ( TLS handshake success) and the request
	actualResponse := readResponseBody(response)

	// Check if the response body is as expected for the second request ( TLS handshake success) and the request
	if actualResponse != expectedResponse {
		t.Fatalf("Expected response body %s, got %s", expectedResponse, actualResponse)
	}

}

// Helper function for reading the response body
func readResponseBody(response *http.Response) string {
	body, _ := io.ReadAll(response.Body)   // read the response body
	defer closeResponseBody(response.Body) // close the response body
	return string(body)                    // return the response body as a string
}

// TestSendRequest_Success tests the sendRequest function with a successful HTTP request
func TestSendRequest_Success(t *testing.T) {
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
		URL:      server.URL,     // URL of the mock server
		Method:   http.MethodGet, // Default to GET method
		Payload:  nil,            // No payload
		Username: "",             // No username
		Password: "",             // No password
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

	// Check if the response body is as expected
	actualResponse := readResponseBody(response)

	// Check if the response body is as expected
	if actualResponse != expectedResponse {
		t.Fatalf("Expected response body %s, got %s", expectedResponse, actualResponse)
	}
}

// TestSendRequest_TlsErrorRetry tests the sendRequest function with TLS error retry logic
func TestSendRequest_TlsErrorRetry(t *testing.T) {
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
		URL:      server.URL,     // URL of the mock server
		Method:   http.MethodGet, // Default to GET method
		Payload:  nil,            // No payload
		Username: "",             // No username
		Password: "",             // No password
	}

	// Create an HTTP client with a transport that injects a TLS error
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

	// Force a TLS error
	client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
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

	// Check if the response body is as expected
	actualResponse := readResponseBody(response)

	// Check if the response body is as expected
	if actualResponse != expectedResponse {
		t.Fatalf("Expected response body %s, got %s", expectedResponse, actualResponse)
	} else {
		t.Logf("Expected response body %s, got %s", expectedResponse, actualResponse)
	}

}

// TestSendRequest_TlsErrorRetrySendRequest tests the sendRequest function with TLS error retry logic
func TestSendRequest_TlsErrorRetrySendRequest(t *testing.T) {
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
		URL:      server.URL,     // URL of the mock server
		Method:   http.MethodGet, // Default to GET method
		Payload:  nil,            // No payload
		Username: "",             // No username
		Password: "",             // No password
	}

	// Create an HTTP client with a transport that injects a TLS error
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

	// Force a TLS error
	client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
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

	// Check if the response body is as expected
	actualResponse := readResponseBody(response)

	// Check if the response body is as expected
	if actualResponse != expectedResponse {
		t.Fatalf("Expected response body %s, got %s", expectedResponse, actualResponse)
	}
}

// TestSendRequest_TlsErrorRetrySendRequestWithReplaceProtocolToHTTP tests the sendRequest function with TLS error retry logic
func TestSendRequest_TlsErrorRetrySendRequestWithReplaceProtocolToHTTP(t *testing.T) {
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

	// Create an HTTP client with a transport that injects a TLS error
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{},
		},
	}

	// replace protocol to http
	config.URL = strings.Replace(config.URL, "https", "http", 1)

	// Create an HTTP request
	request, err := createRequest(context.Background(), config.Method, config.URL, nil, config.Username, config.Password)
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	// Force a TLS error
	client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
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

	// Check if the response body is as expected
	actualResponse := readResponseBody(response)

	// Check if the response body is as expected
	if actualResponse != expectedResponse {
		t.Fatalf("Expected response body %s, got %s", expectedResponse, actualResponse)
	}
}

// Helper function to create a sample request configuration
func createSampleConfig(url string) requestConfig {
	return requestConfig{
		Method:   http.MethodGet, // Default to GET method
		URL:      url,            // URL of the mock server
		Username: "test",         // Username for BasicAuth
		Password: "password",     // Password for BasicAuth
	}
}

// Helper function for assertError to check if an error is not nil
func assertError(t *testing.T, err error, message string) {
	t.Helper() // This line is needed to tell the test suite that this method is a helper method
	if err == nil {
		t.Error(message + ", got nil")
	}
}

// Helper function for assertErrorContains to check if an error contains a substring
func assertErrorContains(t *testing.T, err error, substring string) {

	t.Helper() // This line is needed to tell the test suite that this method is a helper method

	// Check if there is an error and if the error contains the substring
	if err == nil || !strings.Contains(err.Error(), substring) {
		t.Errorf("Expected an error containing '%s', got %v", substring, err)
	}
}

// Helper function to set up a mock server with a status code and response body
func setupMockServer(statusCode int, responseBody string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		_, err := w.Write([]byte(responseBody))
		if err != nil {
			return
		}
	}))
}

// Helper functions for assertions to check if an error is nil for a successful HTTP request
func assertNoError(t *testing.T, err error) {
	t.Helper() // This line is needed to tell the test suite that this method is a helper method
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

// assertNotNil function for assertions to check if a value is not nil for a successful HTTP request
func assertNotNil(t *testing.T, v interface{}) {
	t.Helper() // This line is needed to tell the test suite that this method is a helper method
	if v == nil {
		t.Error("Expected a non-nil value")
	}
}

// TestMakeRequest_Successful tests the makeRequest function with a successful HTTP request
func TestMakeRequest_Successful(t *testing.T) {

	// Create a mock server with a 200 status code (OK) and a success response body
	server := setupMockServer(http.StatusOK, `{"status": "success"}`)

	// Close the mock server
	defer server.Close()

	// Create a request configuration with the mock server URL
	config := createSampleConfig(server.URL)

	// Attempt to make a request with a 200 status code (OK) and a success response body
	response, err := makeRequest(context.Background(), config)

	assertNoError(t, err)     // Check if there is no error for a successful HTTP request
	assertNotNil(t, response) // Check if the response is not nil for a successful HTTP request
}

// TestMakeRequest_InvalidURL tests the makeRequest function is tested for an error for an invalid URL
func TestMakeRequest_InvalidURL(t *testing.T) {

	// Create a request configuration with an invalid URL
	config := createSampleConfig("invalid-url")

	// Attempt to make a request with an invalid URL
	_, err := makeRequest(context.Background(), config)

	// Check if there is an error for an invalid URL
	assertErrorContains(t, err, "makeRequest: invalid URL")
}

// TestMakeRequest_Non2xxStatusCode tests the makeRequest function is tested for an error for a non-2xx status code
func TestMakeRequest_Non2xxStatusCode(t *testing.T) {

	// Create a mock server with a 404 status code (Not Found) and an error response body
	server := setupMockServer(http.StatusNotFound, `{"error": "not found"}`)

	// Close the mock server
	defer server.Close()

	// Create a request configuration with the mock server URL
	config := createSampleConfig(server.URL)

	// Attempt to make a request with a 404 status code (Not Found) and an error response body
	_, err := makeRequest(context.Background(), config)

	// Check if there is an error for a non-2xx status code and an error response body
	assertError(t, err, "Expected an error for non-2xx status code")
}

// TestMakeRequest_Non2xxStatusCode_EmptyResponseBody tests the makeRequest function is tested for an error for a non-2xx
func TestMakeRequest_Non2xxStatusCode_EmptyResponseBody(t *testing.T) {

	// Create a mock server with a 404 status code (Not Found) and an empty response body
	server := setupMockServer(http.StatusNotFound, "")

	// Close the mock server
	defer server.Close()

	// Create a request configuration with the mock server URL
	config := createSampleConfig(server.URL)

	// Attempt to make a request with a 404 status code (Not Found) and an empty response body
	_, err := makeRequest(context.Background(), config)

	// Check if there is an error for a non-2xx status code and an empty response body
	assertError(t, err, "Expected an error for non-2xx status code")
}

// TestMakeRequest_Non2xxStatusCode_LargeResponseBody tests the makeRequest function is tested for an error for a non-2xx
func TestMakeRequest_Non2xxStatusCode_LargeResponseBody(t *testing.T) {

	// Create a mock server with a 404 status code (Not Found) and a large response body
	server := setupMockServer(http.StatusNotFound, string(make([]byte, 1024*1024))) // 1 MB

	// Close the mock server
	defer server.Close()

	// Create a request configuration with the mock server URL
	config := createSampleConfig(server.URL)

	// Attempt to make a request with a 404 status code (Not Found) and a large response body
	_, err := makeRequest(context.Background(), config)

	// Check if there is an error for a non-2xx status code and a large response body
	assertError(t, err, "Expected an error for non-2xx status code")
}

// TestMakeRequest_InvalidMethod tests the makeRequest function is tested for an error Invalid HTTP method
func TestMakeRequest_InvalidMethod(t *testing.T) {

	// Create a request configuration with an valid URL
	config := createSampleConfig("http://example.com")

	// Set the HTTP method to an invalid value
	config.Method = "INVALID"

	// Attempt to make a request with an invalid HTTP method
	_, err := makeRequest(context.Background(), config)

	// Check if there is an error for an invalid HTTP method
	assertError(t, err, "Expected an error for invalid HTTP method")
}

// TestMakeRequest_InvalidURLAndMethod tests the makeRequest function is tested for an error parsing URL and HTTP method
func TestMakeRequest_InvalidURLAndMethod(t *testing.T) {

	// Create a request configuration with an invalid URL
	config := createSampleConfig("invalid-url")

	// Set the HTTP method to an invalid value
	config.Method = "INVALID"

	// Attempt to make a request with an invalid URL and HTTP method
	_, err := makeRequest(context.Background(), config)

	// Check if there is an error for an invalid URL and HTTP method
	assertError(t, err, "Expected an error for invalid URL and HTTP method")
}

// TestMakeRequest_ErrorParsingURL tests the makeRequest function is tested for an error parsing URL
func TestMakeRequest_ErrorParsingURL(t *testing.T) {

	// Create a request configuration with an invalid URL
	config := createSampleConfig(":invalid-url")

	// Attempt to make a request with an invalid URL
	_, err := makeRequest(context.Background(), config)

	// Check if there is an error for an error parsing URL
	assertError(t, err, "Expected an error for error parsing URL")
}

// TestMakeRequest_ErrorParsingURLAndMethod tests the makeRequest function is tested for an error parsing URL and HTTP method
// for the test case defined in the test array of struct below it.
func TestMakeRequest_ErrorParsingURLAndMethod(t *testing.T) {

	// Create a request configuration with an invalid URL and HTTP method
	config := createSampleConfig(":invalid-url")

	// Set the HTTP method to an invalid value
	config.Method = "INVALID"

	// Attempt to make a request with an invalid URL and HTTP method
	_, err := makeRequest(context.Background(), config)

	// Check if there is an error for an invalid URL and HTTP method
	assertError(t, err, "Expected an error for invalid URL and HTTP method")
}

// TestHandleHTTPError_BadRequest tests the handleHTTPError function is nil for a 400 status code (Bad Request)
func TestHandleHTTPError_BadRequest(t *testing.T) {

	// Create a response with a 400 status code (Bad Request)
	errorResponse := &http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       io.NopCloser(strings.NewReader(`{"error": "bad request"}`)),
	}

	// Attempt to handle the error
	err := handleHTTPError(errorResponse)

	// Check if there is no error for a 400 status code (Bad Request)
	assertErrorContains(t, err, "HTTP error")
}

// TestHandleHTTPError_NilResponse tests the handleHTTPError function is nil for a nil response
func TestHandleHTTPError_NilResponse(t *testing.T) {

	// Attempt to handle the error
	err := handleHTTPError(nil)

	// Check if there is an error for a nil response
	assertError(t, err, "Expected an error for nil response")
}

// TestHandleHTTPError_EmptyResponseBody tests the handleHTTPError function with various inputs and expected outputs
// for each test case defined in the tests array of struct below it.
// This test case is for an empty response body.
func TestHandleHTTPError_EmptyResponseBody(t *testing.T) {

	// Create a response with an empty body and a 400 status code (Bad Request)
	errorResponse := &http.Response{
		StatusCode: http.StatusBadRequest,
	}

	// Attempt to handle the error
	err := handleHTTPError(errorResponse)

	// Check if there is an error for an empty response body
	assertError(t, err, "Expected an error for empty response body")
}
