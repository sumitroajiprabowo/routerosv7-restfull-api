/*
Package routerosv7_restfull_api provides functionality for working with RouterOSv7 via RESTful API.
This is unittest for protocol.go
*/
package routerosv7_restfull_api

import (
	"errors"
	"net"
	"testing"
	"time"
)

/*
mockDialer is a mock dialer for testing purposes only that returns an error
on dial and close calls to simulate a connection failure to a host on a port that is not available.
(e.g.localhost:9999) or a connection that is closed by the remote host.
*/
type mockErrorConn struct{}

func (m *mockErrorConn) Read([]byte) (n int, err error) {
	return 0, nil
}

func (m *mockErrorConn) Write([]byte) (n int, err error) {
	return 0, nil
}

func (m *mockErrorConn) Close() error {
	return errors.New("mocked error on close")
}

func (m *mockErrorConn) LocalAddr() net.Addr {
	return nil
}

func (m *mockErrorConn) RemoteAddr() net.Addr {
	return nil
}

func (m *mockErrorConn) SetDeadline(time.Time) error {
	return nil
}

func (m *mockErrorConn) SetReadDeadline(time.Time) error {
	return nil
}

func (m *mockErrorConn) SetWriteDeadline(time.Time) error {
	return nil
}

/*
TestIsHostAvailableOnPort_NotAvailable tests the isHostAvailableOnPort function.
It is not possible to test the actual connection to the host, if the host is not available.
*/
func TestIsHostAvailableOnPort_NotAvailable(t *testing.T) {
	// Set the dialerInstance to use the mockDialer for testing

	defer func() {}()

	// Use a non-existent port for testing
	available := isHostAvailableOnPort("localhost", "9999")
	if available {
		t.Error("Expected host to be not available, got true")
	}
}

/*
TestIsHostAvailableOnPort_Available tests the isHostAvailableOnPort function.
It is not possible to test the actual connection to the host, if the host is not available.
*/
func TestIsHostAvailableOnPort_Available(t *testing.T) {
	// Create a listener to simulate a server on localhost:8080
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		t.Error("Failed to create a listener for coverage")
		return
	}

	// Close the listener when the test is done
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			t.Error("Failed to close the listener for coverage")
		}
	}(listener)

	// Set the dialerInstance to use the defaultDialer for testing
	defer func() {}()

	// Use an existing port for testing
	available := isHostAvailableOnPort("localhost", "8080")
	if !available {
		t.Error("Expected host to be available, got false")
	}
}

/*
TestShouldRetryRequest_TLSHandshakeFailure_HTTPS tests the shouldRetryTlsErrorRequest function.
It is not possible to test the actual connection to the host, if https is not available.
*/
func TestShouldRetryRequest_TLSHandshakeFailure_HTTP(t *testing.T) {

	// Test retrying the request
	err := errors.New("tls: handshake failure")

	// Test retrying the request
	retry := shouldRetryTlsErrorRequest(err, httpProtocol)
	if retry {
		t.Error("Expected no retry for TLS handshake failure with HTTP, got true")
	}
}

/*
TestShouldRetryRequest_TLSHandshakeFailure_HTTPS tests the shouldRetryTlsErrorRequest function.
It is not possible to test the actual connection to the host, if https is not available.
This is because the URL is not parsed in the function.
*/
func TestDetermineProtocol_HTTP(t *testing.T) {

	// Assuming the host is not actually available on port 80 in the test environment
	protocol := determineProtocol("example.com:80")

	// Test the actual connection to the host
	if protocol != httpProtocol {
		t.Errorf("Expected HTTP protocol, got %s", protocol)
	}
}

/*
TestDetermineProtocol_HTTPS tests the determineProtocol function.
It is not possible to test the actual connection to the host, if https is not available.
This is because the URL is not parsed in the function.
*/
func TestDetermineProtocol_HTTPS(t *testing.T) {

	// Assuming the host is not actually available on port 443 in the test environment
	protocol := determineProtocol("example.com")

	// Test the actual connection to the host
	if protocol != httpsProtocol {
		t.Errorf("Expected HTTPS protocol, got %s", protocol)
	}
}

/*
TestCloseResponseBody tests the closeResponseBody function.
It is not possible to test the actual closing of the response body,
if the response body is nil.
*/
func TestCloseConnection(t *testing.T) {

	// Test closing the response body
	mockConn := &mockErrorConn{}

	// Test closing the response body
	closeConnection(mockConn)
}

/*
TestDetermineProtocolFromURL_HTTP tests the determineProtocolFromURL function.
It is not possible to test the actual connection to the host, if http is not available.
This is because the URL is not parsed in the function.
If the host is not available on port 80, the function will return the HTTP protocol.
*/
func TestDetermineProtocolFromURL_HTTP(t *testing.T) {
	// Assuming the host is not actually available on port 80 in the test environment
	protocol := determineProtocolFromURL("http://example.com")

	// Test the actual connection to the host
	if protocol != httpProtocol {
		t.Errorf("Expected HTTP protocol, got %s", protocol)
	}
}

/*
TestDetermineProtocolFromURL_HTTPS tests the determineProtocolFromURL function.
It is not possible to test the actual connection to the host, if https is not available.
This is because the URL is not parsed in the function.
If the host is not available on port 443, the function will return the HTTPS protocol.
*/
func TestDetermineProtocolFromURL_HTTPS(t *testing.T) {
	// Assuming the host is not actually available on port 443 in the test environment
	protocol := determineProtocolFromURL("https://example.com")

	// Test the actual connection to the host
	if protocol != httpsProtocol {
		t.Errorf("Expected HTTPS protocol, got %s", protocol)
	}
}

/*
TestReplaceProtocol tests the replaceProtocol function.
It is not possible to test the actual replacement of the protocol in the URL,
if https is not available. This is because the URL is not parsed in the function.
*/
func TestReplaceProtocol(t *testing.T) {

	// Test replacing the protocol from HTTP to HTTPS
	newURL := replaceProtocol("http://example.com", httpProtocol, httpsProtocol)

	// Test replacing the protocol from HTTPS to HTTP (not possible to test)
	if newURL != "https://example.com" {
		t.Errorf("Expected URL with HTTPS protocol, got %s", newURL)
	}
}
