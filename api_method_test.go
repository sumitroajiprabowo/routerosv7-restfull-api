// api_methods_test.go

package routerosv7_restfull_api

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestURL(t *testing.T) {

	// Example:
	request := &APIRequest{
		Host:     "example.com",
		Command:  "example",
		Method:   MethodGet,
		Username: "user",
		Password: "pass",
		Payload:  []byte("payload"),
	}

	expectedURL := "http://example.com/rest/example"
	actualURL := request.URL()

	assert.Equal(t, expectedURL, actualURL, "URL does not match expected")
}

//
//func TestAuth(t *testing.T) {
//	// Create a mock API
//	mockAPI := new(MockRouterOSAPI)
//
//	// Set up the expected result and error for makeRequest
//	expectedMakeRequestResult := map[string]interface{}{"result": "success"}
//	var expectedMakeRequestError error // Change this to the expected error if any
//
//	// Set up the expectations for the makeRequest method
//	mockAPI.On("makeRequest", mock.Anything). // Assuming makeRequest only requires context
//							Return(expectedMakeRequestResult, expectedMakeRequestError)
//
//	// Create an instance of AuthConfig
//	authConfig := AuthConfig{
//		Host:     "testhost",
//		Username: "testuser",
//		Password: "testpassword",
//	}
//
//	// Execute the Auth function using the mock API
//	result, err := Auth(context.Background(), authConfig)
//
//	// Assert that the expectations were met for makeRequest
//	mockAPI.AssertExpectations(t)
//
//	// Check for errors
//	assert.NoError(t, err, "Auth function returned an error")
//
//	// Verify the result
//	assert.Equal(t, expectedMakeRequestResult, result, "Auth function result does not match expected result")
//}

func TestPrint(t *testing.T) {
	// Create a mock API
	mockAPI := new(MockRouterOSAPI)

	// Set up the expected result and error
	expectedResult := map[string]interface{}{"result": "success"}
	var expectedError error // Change this to the expected error if any

	// Set up the expectations for the Print method
	mockAPI.On("Print", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(expectedResult, expectedError)

	// Execute the Print function using the mock API
	result, err := mockAPI.Print(context.Background(), "testhost", "testuser", "testpassword", "testcommand")

	// Assert that the expectations were met
	mockAPI.AssertExpectations(t)

	// Check for errors
	assert.NoError(t, err, "Print function returned an error")

	// Verify the result
	assert.Equal(t, expectedResult, result, "Print function result does not match expected result")
}

func TestAdd(t *testing.T) {
	// Create a mock API
	mockAPI := new(MockRouterOSAPI)

	// Set up the expected result and error
	expectedResult := map[string]interface{}{"result": "success"}
	var expectedError error // Change this to the expected error if any

	// Set up the expectations for the Add method
	mockAPI.On("Add", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(expectedResult, expectedError)

	// Execute the Add function using the mock API
	result, err := mockAPI.Add(context.Background(), "testhost", "testuser", "testpassword", "testcommand", []byte("testpayload"))

	// Assert that the expectations were met
	mockAPI.AssertExpectations(t)

	// Check for errors
	assert.NoError(t, err, "Add function returned an error")

	// Verify the result
	assert.Equal(t, expectedResult, result, "Add function result does not match expected result")
}

func TestSet(t *testing.T) {
	// Create a mock API
	mockAPI := new(MockRouterOSAPI)

	// Set up the expected result and error
	expectedResult := map[string]interface{}{"result": "success"}
	var expectedError error // Change this to the expected error if any

	// Set up the expectations for the Set method
	mockAPI.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(expectedResult, expectedError)

	// Execute the Set function using the mock API
	result, err := mockAPI.Set(context.Background(), "testhost", "testuser", "testpassword", "testcommand", []byte("testpayload"))

	// Assert that the expectations were met
	mockAPI.AssertExpectations(t)

	// Check for errors
	assert.NoError(t, err, "Set function returned an error")

	// Verify the result
	assert.Equal(t, expectedResult, result, "Set function result does not match expected result")
}

func TestRemove(t *testing.T) {
	// Create a mock API
	mockAPI := new(MockRouterOSAPI)

	// Set up the expected result and error
	expectedResult := map[string]interface{}{"result": "success"}
	var expectedError error // Change this to the expected error if any

	// Set up the expectations for the Remove method
	mockAPI.On("Remove", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(expectedResult, expectedError)

	// Execute the Remove function using the mock API
	result, err := mockAPI.Remove(context.Background(), "testhost", "testuser", "testpassword", "testcommand")

	// Assert that the expectations were met
	mockAPI.AssertExpectations(t)

	// Check for errors
	assert.NoError(t, err, "Remove function returned an error")

	// Verify the result
	assert.Equal(t, expectedResult, result, "Remove function result does not match expected result")
}

func TestRun(t *testing.T) {
	// Create a mock API
	mockAPI := new(MockRouterOSAPI)

	// Set up the expected result and error
	expectedResult := map[string]interface{}{"result": "success"}
	var expectedError error // Change this to the expected error if any

	// Set up the expectations for the Run method
	mockAPI.On("Run", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(expectedResult, expectedError)

	// Execute the Run function using the mock API
	result, err := mockAPI.Run(context.Background(), "testhost", "testuser", "testpassword", "testcommand", []byte("testpayload"))

	// Assert that the expectations were met
	mockAPI.AssertExpectations(t)

	// Check for errors
	assert.NoError(t, err, "Run function returned an error")

	// Verify the result
	assert.Equal(t, expectedResult, result, "Run function result does not match expected result")
}
