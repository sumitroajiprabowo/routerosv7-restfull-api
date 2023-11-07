package routerosv7_restfull_api

import (
	"github.com/jarcoal/httpmock"
	"testing"
)

func TestDetermineProtocol_HTTP(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Mock server HTTP yang selalu mengembalikan respons HTTP 404
	httpmock.RegisterResponder("GET", "http://example.com", httpmock.NewStringResponder(404, ""))

	protocol := determineProtocol("http://example.com")

	if protocol != "http" {
		t.Errorf("Expected protocol to be 'http', but got '%s'", protocol)
	}
}
