// api_methods_test.go
package routerosv7_restfull_api

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
