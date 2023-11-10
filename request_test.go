package routerosv7_restfull_api

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMakeRequest(t *testing.T) {
	// Mocking a server for testing
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "success"}`))
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

	// Simulate an error during the request
	config.URL = "invalid-url" // This will cause a request error
	_, err = makeRequest(context.Background(), config)
	if err == nil {
		t.Error("Expected an error, got nil")
	}

	// Testing non-2xx status code handling
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "not found"}`))
	}))
	defer server.Close()

	config.URL = server.URL
	_, err = makeRequest(context.Background(), config)
	if err == nil {
		t.Error("Expected an error for non-2xx status code, got nil")
	}

	// Testing handleHTTPError function
	errorResponse := &http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       ioutil.NopCloser(strings.NewReader(`{"error": "bad request"}`)),
	}

	err = handleHTTPError(errorResponse)
	if err == nil {
		t.Error("Expected an error, got nil")
	}
}

func TestMakeRequest_InvalidHTTPMethod(t *testing.T) {
	// Mocking a server for testing
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "success"}`))
	}))
	defer server.Close()

	// Sample request configuration with an invalid HTTP method
	config := requestConfig{
		Method:   "INVALID_METHOD",
		URL:      server.URL,
		Username: "test",
		Password: "password",
	}

	// Testing an error for an invalid HTTP method
	_, err := makeRequest(context.Background(), config)
	if err == nil {
		t.Error("Expected an error for invalid HTTP method, got nil")
	} else if !strings.Contains(err.Error(), "makeRequest: invalid HTTP method") {
		t.Errorf("Expected an error containing 'makeRequest: invalid HTTP method', got %v", err)
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

}

func TestHandleHTTPErrorReadErrorNonNilBody(t *testing.T) {
	// Test case 9: Read error on non-nil HTTP response body
	statusCode9 := http.StatusNotFound
	response9 := &http.Response{
		Status:     fmt.Sprintf("%d %s", statusCode9, http.StatusText(statusCode9)),
		StatusCode: statusCode9,
		Body:       &errorReader{returnError: true},
	}

	err9 := handleHTTPError(response9)

	// Check if the error message contains the expected HTTP status code and status text
	expectedError9 := fmt.Sprintf("HTTP error: %d %s, Response body: ", statusCode9, http.StatusText(statusCode9))
	if err9 == nil || !strings.Contains(err9.Error(), expectedError9) {
		t.Errorf("handleHTTPError did not return the expected error for Read error on non-nil HTTP response body. Got: %v, Expected: %s", err9, expectedError9)
	}
}

func TestHandleHTTPErrorEmptyBody(t *testing.T) {
	// Create a sample HTTP response with a non-2xx status code and an empty response body
	statusCode := http.StatusNotFound
	response := httptest.NewRecorder()
	response.WriteHeader(statusCode)

	// Call handleHTTPError with the sample response
	err := handleHTTPError(response.Result())

	// Check if the error matches the expected format
	expectedError := fmt.Sprintf("HTTP error: %d %s, Response body: ", statusCode, http.StatusText(statusCode))
	if err == nil || err.Error() != expectedError {
		t.Errorf("handleHTTPError did not return the expected error. Got: %v, Expected: %s", err, expectedError)
	}
}

func TestHandleHTTPErrorReadError(t *testing.T) {
	// Create a sample HTTP response with a non-2xx status code and a response body
	statusCode := http.StatusNotFound
	response := httptest.NewRecorder()
	response.WriteHeader(statusCode)
	response.Body.WriteString("Error response body")

	// Create a custom response body that returns an error on read
	readErrorBody := &errorReader{}

	// Create a bytes.Buffer and set its Reader to our custom body
	buffer := &bytes.Buffer{}
	buffer.ReadFrom(readErrorBody)

	// Set the response body to our buffer
	response.Body = buffer

	// Call handleHTTPError with the sample response
	err := handleHTTPError(response.Result())

	// Check if the error matches the expected format
	expectedError := fmt.Sprintf("HTTP error: %d %s, Response body: ", statusCode, http.StatusText(statusCode))
	if err == nil || err.Error() != expectedError {
		t.Errorf("handleHTTPError did not return the expected error. Got: %v, Expected: %s", err, expectedError)
	}
}

// errorReader is a custom io.Reader that returns an error on read
type errorReader struct {
	returnError bool
	statusCode  int
}

func (er *errorReader) Read(_ []byte) (n int, err error) {
	er.statusCode = http.StatusNotFound // You can customize this based on your test case
	return 0, fmt.Errorf("HTTP error: %d %s, Response body: ", er.statusCode, http.StatusText(er.statusCode))
}

func (er *errorReader) Close() error {
	return nil
}

func TestHandleHTTPErrorReadErrorNilBody(t *testing.T) {
	// Create a sample HTTP response with a nil response body
	statusCode := http.StatusNotFound
	response := httptest.NewRecorder()
	response.WriteHeader(statusCode)

	// Set the response body to nil
	response.Body = nil

	// Call handleHTTPError with the sample response
	err := handleHTTPError(response.Result())

	// Check if the error matches the expected format
	expectedError := fmt.Sprintf("HTTP error: %d %s, Response body: ", statusCode, http.StatusText(statusCode))
	if err == nil || err.Error() != expectedError {
		t.Errorf("handleHTTPError did not return the expected error for Read error on nil HTTP response body. Got: %v, Expected: %s", err, expectedError)
	}
}

func TestHandleHTTPErrorReadErrorNotNilBody(t *testing.T) {
	// Create a sample HTTP response with a non-nil response body
	statusCode := http.StatusNotFound
	response := httptest.NewRecorder()
	response.WriteHeader(statusCode)
	response.Body.WriteString("Error response body")

	// Create a custom response body that returns an error on read
	readErrorBody := &errorReader{}

	// Create a bytes.Buffer and set its Reader to our custom body
	buffer := &bytes.Buffer{}
	buffer.ReadFrom(readErrorBody)

	// Set the response body to our buffer
	response.Body = buffer

	// Call handleHTTPError with the sample response
	err := handleHTTPError(response.Result())

	// Check if the error matches the expected format
	expectedError := fmt.Sprintf("HTTP error: %d %s, Response body: ", statusCode, http.StatusText(statusCode))
	if err == nil || err.Error() != expectedError {
		t.Errorf("handleHTTPError did not return the expected error for Read error on not nil HTTP response body. Got: %v, Expected: %s", err, expectedError)
	}
}

func TestHandleHTTPErrorNilResponse(t *testing.T) {
	var response *http.Response

	// Call handleHTTPError with nil response
	err := handleHTTPError(response)

	// Check if the error matches the expected format
	expectedError := "nil HTTP response"
	if err == nil || err.Error() != expectedError {
		t.Errorf("handleHTTPError did not return the expected error. Got: %v, Expected: %s", err, expectedError)
	}
}

func TestHandleHTTPErrorNilResponseBody(t *testing.T) {
	// Create a sample HTTP response with a non-2xx status code and a nil response body
	statusCode := http.StatusNotFound
	response := &http.Response{
		Status:     fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode)),
		StatusCode: statusCode,
		Body:       nil,
	}

	// Call handleHTTPError with the sample response
	err := handleHTTPError(response)

	// Check if the error matches the expected format
	expectedError := "nil HTTP response body"
	if err == nil || err.Error() != expectedError {
		t.Errorf("handleHTTPError did not return the expected error. Got: %v, Expected: %s", err, expectedError)
	}
}

//func TestCreateRequest_Error(t *testing.T) {
//	config := requestConfig{
//		Method:   http.MethodGet,
//		URL:      "http://invalid-url", // Provide a valid URL with a scheme
//		Username: "test",
//		Password: "password",
//	}
//
//	// Test case for createRequest error
//	request, err := createRequest(context.Background(), config.Method, config.URL, nil, config.Username, config.Password)
//	if err == nil {
//		t.Error("Expected an error for invalid URL, got nil")
//	} else {
//		// Print the actual error for debugging
//		fmt.Printf("Actual Error: %v\n", err)
//	}
//
//	// Additional check to see if the returned request is nil
//	if request != nil {
//		t.Error("Expected a nil request for invalid URL")
//	} else {
//		fmt.Println("Request is nil as expected.")
//	}
//}

//func TestMakeRequest_Retry(t *testing.T) {
//	// Mocking a server for testing
//	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		w.WriteHeader(http.StatusInternalServerError)
//	}))
//	defer server.Close()
//
//	// Sample request configuration
//	config := requestConfig{
//		Method:   http.MethodGet,
//		URL:      server.URL,
//		Username: "test",
//		Password: "password",
//	}
//
//	// Test retry scenario
//	response, err := makeRequest(context.Background(), config)
//	if err != nil {
//		t.Errorf("Expected no error, got %v", err)
//	}
//
//	// Ensure the response is not nil
//	if response == nil {
//		t.Error("Expected a non-nil response")
//	}
//}

//func TestCreateRequest_InvalidURL(t *testing.T) {
//	config := requestConfig{
//		Method:   http.MethodGet,
//		URL:      "invalid-url",
//		Username: "test",
//		Password: "password",
//	}
//
//	request, err := createRequest(context.Background(), config.Method, config.URL, nil, config.Username, config.Password)
//	if err == nil {
//		t.Error("Expected an error for invalid URL, got nil")
//	}
//
//	// Additional check to see if the returned request is nil
//	if request != nil {
//		t.Error("Expected a nil request for invalid URL")
//	}
//}

// Test case for handling non-2xx status codes
func TestMakeRequest_Non2xxStatusCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "not found"}`))
	}))
	defer server.Close()

	config := requestConfig{
		Method:   http.MethodGet,
		URL:      server.URL,
		Username: "test",
		Password: "password",
	}

	_, err := makeRequest(context.Background(), config)
	if err == nil {
		t.Error("Expected an error for non-2xx status code, got nil")
	}
}

// Test case for handling invalid URLs
func TestMakeRequest_InvalidURL(t *testing.T) {
	config := requestConfig{
		Method:   http.MethodGet,
		URL:      "invalid-url",
		Username: "test",
		Password: "password",
	}

	_, err := makeRequest(context.Background(), config)
	if err == nil {
		t.Error("Expected an error for invalid URL, got nil")
	}
}

// Test case for handling non-JSON response body
func TestMakeRequest_NonJSONResponseBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("plain text response"))
	}))
	defer server.Close()

	config := requestConfig{
		Method:   http.MethodGet,
		URL:      server.URL,
		Username: "test",
		Password: "password",
	}

	_, err := makeRequest(context.Background(), config)
	if err == nil {
		t.Error("Expected an error for non-JSON response body, got nil")
	}
}

func TestIsHostAvailableOnPort(t *testing.T) {
	// Valid host and port
	available := isHostAvailableOnPort("example.com", "80")
	if !available {
		t.Error("Expected host to be available on port 80")
	}

	// Invalid host and port
	available = isHostAvailableOnPort("invalidhost", "8080")
	if available {
		t.Error("Expected host to be unavailable on port 8080")
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

func TestCreateRequestBody(t *testing.T) {
	// Non-empty payload
	payload := []byte(`{"key": "value"}`)
	body := createRequestBody(payload)
	if body == nil {
		t.Error("Expected non-nil body for non-empty payload")
	}

	// Empty payload
	emptyPayload := []byte{}
	emptyBody := createRequestBody(emptyPayload)
	if emptyBody != nil {
		t.Error("Expected nil body for empty payload")
	}
}

func TestCloseResponseBody(t *testing.T) {
	// Mock a response body with an error on close
	errorBody := &mockErrorReaderCloser{}
	closeResponseBody(errorBody) // This should log the error, you can capture logs and check
}

type mockErrorReaderCloser struct{}

func (m *mockErrorReaderCloser) Read(p []byte) (n int, err error) {
	return 0, io.EOF
}

func (m *mockErrorReaderCloser) Close() error {
	return errors.New("mocked close error")
}

func TestMakeRequestRetrySuccess(t *testing.T) {
	// Simulate an initial request failure
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	config := requestConfig{
		Method:   http.MethodGet,
		URL:      server.URL,
		Username: "test",
		Password: "password",
	}

	// Initial request failure
	_, err := makeRequest(context.Background(), config)
	if err == nil {
		t.Error("Expected an error, got nil")
	}

	// Retry with success
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "success"}`))
	}))
	defer server.Close()

	config.URL = server.URL
	response, err := makeRequest(context.Background(), config)
	if err != nil {
		t.Errorf("Expected no error on retry, got %v", err)
	}

	// Ensure the response is not nil
	if response == nil {
		t.Error("Expected a non-nil response on retry")
	}
}

func TestDecodeEmptyJSONBody(t *testing.T) {
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
}

type emptyReadCloser struct{}

func (erc *emptyReadCloser) Read(p []byte) (n int, err error) {
	return 0, io.EOF
}

func (erc *emptyReadCloser) Close() error {
	return nil
}
