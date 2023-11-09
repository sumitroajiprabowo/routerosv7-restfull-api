package routerosv7_restfull_api

import (
	"testing"
	"time"
)

func TestConstants(t *testing.T) {
	// Test httpProtocol
	if httpProtocol != "http" {
		t.Errorf("Expected httpProtocol to be 'http', got %s", httpProtocol)
	}

	// Test httpsProtocol
	if httpsProtocol != "https" {
		t.Errorf("Expected httpsProtocol to be 'https', got %s", httpsProtocol)
	}

	// Test tlsHandshakeFailure
	if tlsHandshakeFailure != "tls: handshake failure" {
		t.Errorf("Expected tlsHandshakeFailure to be 'tls: handshake failure', got %s", tlsHandshakeFailure)
	}

	// Test pingCount
	if pingCount != 3 {
		t.Errorf("Expected pingCount to be 3, got %d", pingCount)
	}

	// Test pingTimeout
	expectedTimeout := 1000 * time.Millisecond
	if pingTimeout != expectedTimeout {
		t.Errorf("Expected pingTimeout to be %s, got %s", expectedTimeout, pingTimeout)
	}

	// Test pingInterval
	expectedInterval := 100 * time.Millisecond
	if pingInterval != expectedInterval {
		t.Errorf("Expected pingInterval to be %s, got %s", expectedInterval, pingInterval)
	}
}

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
