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

func TestIsValidURL(t *testing.T) {
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidURL(tt.url)
			if result != tt.expected {
				t.Errorf("isValidURL(%s) = %v; want %v", tt.url, result, tt.expected)
			}
		})
	}
}

func mustParseURL(rawURL string) *url.URL {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}
	return parsedURL
}

func TestIsValidHTTPMethod(t *testing.T) {
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidHTTPMethod(tt.method)
			if result != tt.expected {
				t.Errorf("isValidHTTPMethod(%s) = %v; want %v", tt.method, result, tt.expected)
			}
		})
	}
}

func TestParseURL(t *testing.T) {
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseURL(tt.rawURL)

			if (err != nil) != tt.wantErr {
				t.Errorf("parseURL(%s) error = %v, wantErr %v", tt.rawURL, err, tt.wantErr)
				return
			}

			if tt.expected != nil && (result == nil || result.String() != tt.expected.String()) {
				t.Errorf("parseURL(%s) = %v, want %v", tt.rawURL, result, tt.expected)
			}
		})
	}
}

func TestCreateRequestBody(t *testing.T) {
	// Non-empty payload
	payload := []byte(`{"key": "value"}`)
	body := createRequestBody(payload)
	if body == nil {
		t.Error("Expected non-nil body for non-empty payload")
	}

	// Empty payload
	var emptyPayload []byte
	emptyBody := createRequestBody(emptyPayload)
	if emptyBody != nil {
		t.Error("Expected nil body for empty payload")
	}
}

type mockErrorReaderCloser struct{}

func (m *mockErrorReaderCloser) Read(_ []byte) (n int, err error) {
	return 0, errors.New("mocked read error")
}

func (m *mockErrorReaderCloser) Close() error {
	return errors.New("mocked close error")
}

func TestCloseResponseBody(t *testing.T) {
	// Mock a response body with an error on close
	errorBody := &mockErrorReaderCloser{}
	closeResponseBody(errorBody) // This should log the error, you can capture logs and check
}

func TestValidateRequestConfig(t *testing.T) {
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRequestConfig(tt.config)

			if tt.expectedError && err == nil {
				t.Error("validateRequestConfig() did not return the expected error.")
			}

			if !tt.expectedError && err != nil {
				t.Errorf("validateRequestConfig() returned an unexpected error: %v", err)
			}
		})
	}
}

func TestCreateHTTPClient(t *testing.T) {
	// Test for HTTPS protocol
	client := createHTTPClient(httpsProtocol)
	if client.Transport == nil || client.Transport.(*http.Transport).TLSClientConfig == nil {
		t.Error("Expected TLS config to be set for HTTPS protocol")
	}

	// Test for HTTP protocol
	client = createHTTPClient(httpProtocol)
	if client.Transport != nil {
		t.Error("Expected no TLS config for HTTP protocol")
	}
}

func TestHandleHTTPError(t *testing.T) {
	// Test case 1: Non-2xx status code with a response body
	responseBody1 := "Error response body 1"
	statusCode1 := http.StatusNotFound
	response1 := httptest.NewRecorder()
	response1.WriteHeader(statusCode1)
	response1.Body.WriteString(responseBody1)

	// Call handleHTTPError with the first sample response
	err1 := handleHTTPError(response1.Result())

	// Check if the error matches the expected format
	expectedError1 := fmt.Sprintf("HTTP error: %d %s, Response body: %s", statusCode1, http.StatusText(statusCode1), responseBody1)
	if err1 == nil || err1.Error() != expectedError1 {
		t.Errorf("handleHTTPError did not return the expected error. Got: %v, Expected: %s", err1, expectedError1)
	}

	// Test case 2: Non-2xx status code without a response body
	statusCode2 := http.StatusInternalServerError
	response2 := httptest.NewRecorder()
	response2.WriteHeader(statusCode2)

	// Call handleHTTPError with the second sample response
	err2 := handleHTTPError(response2.Result())

	// Check if the error matches the expected format
	expectedError2 := fmt.Sprintf("HTTP error: %d %s, Response body: ", statusCode2, http.StatusText(statusCode2))
	if err2 == nil || err2.Error() != expectedError2 {
		t.Errorf("handleHTTPError did not return the expected error. Got: %v, Expected: %s", err2, expectedError2)
	}

	// Test case 3: 2xx status code with a response body
	response3 := httptest.NewRecorder()
	response3.WriteHeader(http.StatusOK)
	response3.Body.WriteString("Success response body")

	// Call handleHTTPError with the third sample response
	err3 := handleHTTPError(response3.Result())

	// Check if the error is not nil for 2xx status codes
	if err3 == nil {
		t.Errorf("handleHTTPError did not return an error for a 2xx status code. Expected: non-nil error")
	}

	// Test case 4: Non-2xx status code with an empty response body
	statusCode4 := http.StatusBadRequest
	response4 := httptest.NewRecorder()
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
	response5 := httptest.NewRecorder()
	response5.WriteHeader(statusCode5)
	largeBody := "This is a very large response body. " + string(make([]byte, 1024*1024)) // 1 MB
	response5.Body.WriteString(largeBody)

	// Call handleHTTPError with the fifth sample response
	err5 := handleHTTPError(response5.Result())

	// Check if the error matches the expected format
	expectedError5 := fmt.Sprintf("HTTP error: %d %s, Response body: %s", statusCode5, http.StatusText(statusCode5), largeBody)
	if err5 == nil || err5.Error() != expectedError5 {
		t.Errorf("handleHTTPError did not return the expected error. Got: %v, Expected: %s", err5, expectedError5)
	}

	// Test case 6: Nil HTTP response
	var responseNil *http.Response
	err6 := handleHTTPError(responseNil)
	expectedError6 := "nil HTTP response"
	if err6 == nil || err6.Error() != expectedError6 {
		t.Errorf("handleHTTPError did not return the expected error. Got: %v, Expected: %s", err6, expectedError6)
	}

	// Test case 7: Nil HTTP response body
	statusCode7 := http.StatusNotFound
	response7 := &http.Response{
		Status:     fmt.Sprintf("%d %s", statusCode7, http.StatusText(statusCode7)),
		StatusCode: statusCode7,
		Body:       nil,
	}
	err7 := handleHTTPError(response7)
	expectedError7 := "nil HTTP response body"
	if err7 == nil || err7.Error() != expectedError7 {
		t.Errorf("handleHTTPError did not return the expected error. Got: %v, Expected: %s", err7, expectedError7)
	}

	// Test case 8: HTTP response body read error
	statusCode8 := http.StatusNotFound
	response8 := &http.Response{
		Status: fmt.Sprintf("%d %s", statusCode8, http.StatusText(statusCode8)),
		Body:   &mockErrorReaderCloser{},
	}

	err8 := handleHTTPError(response8)
	expectedError8 := "HTTP error: 404 Not Found, Response body: "
	if err8 == nil || err8.Error() != expectedError8 {
		t.Errorf("handleHTTPError did not return the expected error. Got: %v, Expected: %s", err8, expectedError8)
	}

	//Test case 9: HTTP response body close error
	statusCode9 := http.StatusNotFound
	response9 := &http.Response{
		Status: fmt.Sprintf("%d %s", statusCode9, http.StatusText(statusCode9)),
		Body:   &mockErrorReaderCloser{},
	}

	err9 := handleHTTPError(response9)
	expectedError9 := "HTTP error: 404 Not Found, Response body: "
	if err9 == nil || err9.Error() != expectedError9 {
		t.Errorf("Error")
	}

}

type emptyReadCloser struct{}

func (erc *emptyReadCloser) Read(_ []byte) (n int, err error) {
	return 0, io.EOF
}

func (erc *emptyReadCloser) Close() error {
	return nil
}

type stringReadCloser struct {
	io.Reader
}

func (src *stringReadCloser) Close() error {
	return nil
}

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

func TestSetRequestAuth(t *testing.T) {
	// Test case 1: Set BasicAuth with valid username and password
	request := &http.Request{Header: make(http.Header)}
	username := "user"
	password := "pass"
	setRequestAuth(request, username, password)

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

func TestSetRequestContentType(t *testing.T) {
	// Test case 1: Content-Type is set to application/json
	request := &http.Request{Header: make(http.Header)}
	setRequestContentType(request)

	if contentType := request.Header.Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Test case 1: Content-Type not set correctly. Expected application/json, got %s", contentType)
	}
}

func TestNewHTTPRequest(t *testing.T) {
	// Test case 1: Valid GET request without authentication
	ctx := context.Background()
	method := http.MethodGet
	urls := "http://example.com"
	var body io.Reader = nil
	username := ""
	password := ""

	request, err := newHTTPRequest(ctx, method, urls, body, username, password)

	if err != nil {
		t.Errorf("Test case 1: Expected no error, but got an error: %v", err)
	}

	if request == nil {
		t.Error("Test case 1: Expected a non-nil request, but got nil")
	}

	// Test case 2: Valid POST request with authentication
	method = http.MethodPost
	body = strings.NewReader(`{"key": "value"}`)
	username = "user"
	password = "pass"

	request, err = newHTTPRequest(ctx, method, urls, body, username, password)

	if err != nil {
		t.Errorf("Test case 2: Expected no error, but got an error: %v", err)
	}

	if request == nil {
		t.Error("Test case 2: Expected a non-nil request, but got nil")
	}

	// Test case 3: Invalid URL
	urls = ":invalid-url"
	_, err = newHTTPRequest(ctx, method, urls, body, username, password)

	if err == nil {
		t.Error("Test case 3: Expected an error for an invalid URL, but got nil")
	}
}

func TestCreateRequest(t *testing.T) {
	// Test case 1: Valid GET request
	ctx := context.Background()
	method := http.MethodGet
	urls := "http://example.com"
	body := io.Reader(nil)
	username := "user"
	password := "pass"

	request, err := createRequest(ctx, method, urls, body, username, password)

	if err != nil {
		t.Errorf("Test case 1: Expected no error, got %v", err)
	}

	if request == nil {
		t.Error("Test case 1: Expected non-nil request, got nil")
	}

	// Test case 2: Valid POST request
	method = http.MethodPost
	request, err = createRequest(ctx, method, urls, body, username, password)

	if err != nil {
		t.Errorf("Test case 2: Expected no error, got %v", err)
	}

	if request == nil {
		t.Error("Test case 2: Expected non-nil request, got nil")
	}

	// Test case 3: Valid request with body
	method = http.MethodPut
	body = bytes.NewBufferString(`{"key": "value"}`)
	request, err = createRequest(ctx, method, urls, body, username, password)

	if err != nil {
		t.Errorf("Test case 3: Expected no error, got %v", err)
	}

	if request == nil {
		t.Error("Test case 3: Expected non-nil request, got nil")
	}

	// Test case 4: Invalid URL
	urls = ":invalid-url"
	_, err = createRequest(ctx, method, urls, body, username, password)

	if err == nil {
		t.Error("Test case 4: Expected error for invalid URL, got nil")
	}

	// Test case 5: Error parsing URL
	urls = "http://example.com"
	_, err = createRequest(ctx, method, urls, body, username, password)

	if err != nil {
		t.Errorf("Test case 5: Expected no error while creating request, got %v", err)
	}

	// Test case 6: Invalid HTTP method
	method = "INVALID"
	_, err = createRequest(ctx, method, urls, body, username, password)

	if err == nil {
		t.Error("Test case 5: Expected error for invalid HTTP method, got nil")
	}

	// Test case 7: Error creating request
	method = http.MethodGet
	urls = "http://example.com"
	username = "user"
	password = "pass"
	_, err = createRequest(ctx, method, urls, body, username, password)

	if err != nil {
		t.Errorf("Test case 7: Expected no error, got %v", err)
	}

	// Test case 8: Error creating request with invalid URL
	method = http.MethodGet
	urls = ":invalid-url"
	username = "user"
	password = "pass"
	_, err = createRequest(ctx, method, urls, body, username, password)

	if err == nil {
		t.Error("Test case 8: Expected error for invalid URL, got nil")
	}

	// Test case 9: Error creating request with invalid HTTP method
	method = "INVALID"
	urls = "http://example.com"
	username = "user"
	password = "pass"
	_, err = createRequest(ctx, method, urls, body, username, password)

	if err == nil {
		t.Error("Test case 9: Expected error for invalid HTTP method, got nil")
	}

	// Test case 10: Error creating request with invalid HTTP method and invalid URL
	method = "INVALID"
	urls = ":invalid-url"
	username = "user"
	password = "pass"
	_, err = createRequest(ctx, method, urls, body, username, password)

	if err == nil {
		t.Error("Test case 10: Expected error for invalid HTTP method and invalid URL, got nil")
	}

	// Test case 11: Error creating request with empty body
	method = http.MethodGet
	urls = "http://example.com"
	username = "user"
	password = "pass"
	body = nil
	_, err = createRequest(ctx, method, urls, body, username, password)

	if err != nil {
		t.Errorf("Test case 11: Expected no error, got %v", err)
	}

}

func TestRetryTlsErrorRequest(t *testing.T) {
	// Membuat server HTTP sederhana untuk tes
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulasi TLS handshake failure pada request pertama
		if r.URL.Scheme == httpsProtocol {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Simulated TLS handshake failure"))
			return
		}

		// Menanggapi request yang berhasil setelah modifikasi protocol
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Success"))
	}))
	defer server.Close()

	// Konfigurasi permintaan untuk digunakan dalam pengujian
	config := requestConfig{
		URL:      server.URL,
		Method:   http.MethodGet,
		Payload:  nil,
		Username: "",
		Password: "",
	}

	// Membuat klien HTTP palsu dengan transport yang dimodifikasi untuk menggagalkan TLS handshake
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{},
		},
	}

	// Membuat permintaan HTTP palsu
	request, err := http.NewRequest(http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	// Menjalankan fungsi retryTlsErrorRequest untuk pengujian
	response, err := retryTlsErrorRequest(client, request, config)

	// Memeriksa apakah tidak ada kesalahan yang terjadi dan respons diterima
	if err != nil {
		t.Fatalf("RetryTlsErrorRequest failed: %v", err)
	}

	// Memeriksa apakah respons memiliki kode status yang benar setelah perbaikan
	if response.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, response.StatusCode)
	}

	// Memeriksa apakah hasil respons sesuai dengan harapan
	expectedResponse := "Success"
	actualResponse := readResponseBody(response)
	if actualResponse != expectedResponse {
		t.Fatalf("Expected response body %s, got %s", expectedResponse, actualResponse)
	}
}

// Helper function untuk membaca respons body sebagai string
func readResponseBody(response *http.Response) string {
	body, _ := io.ReadAll(response.Body)
	defer closeResponseBody(response.Body)
	return string(body)
}

func TestSendRequest(t *testing.T) {
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Success"))
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
	err = handleHTTPError(errorResponse)

	if err == nil {
		t.Error("Expected an error, got nil")
	}

}
