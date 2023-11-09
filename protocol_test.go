package routerosv7_restfull_api

import (
	"errors"
	"net"
	"testing"
	"time"
)

// dialer is an interface representing the dial function.
type dialer interface {
	dial(network, address string) (net.Conn, error)
}

// defaultDialer is an implementation of the dialer interface using the actual net.Dial function.
type defaultDialer struct{}

// mockDialer is a mock implementation of the dialer interface for testing purposes.
type mockDialer struct {
	called bool
}

var dialerInstance dialer = &defaultDialer{}

func (d *defaultDialer) dial(network, address string) (net.Conn, error) {
	return net.Dial(network, address)
}

type mockErrorConn struct{}

func (d *mockDialer) dial(network, address string) (net.Conn, error) {
	d.called = true
	return nil, errors.New("mocked error")
}

func (m *mockErrorConn) Read(b []byte) (n int, err error) {
	return 0, nil
}

func (m *mockErrorConn) Write(b []byte) (n int, err error) {
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

func (m *mockErrorConn) SetDeadline(t time.Time) error {
	return nil
}

func (m *mockErrorConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (m *mockErrorConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func TestIsHostAvailableOnPort_NotAvailable(t *testing.T) {
	// Set the dialerInstance to use the mockDialer for testing
	dialerInstance = &mockDialer{}
	defer func() { dialerInstance = &defaultDialer{} }()

	// Use a non-existent port for testing
	available := isHostAvailableOnPort("localhost", "9999")
	if available {
		t.Error("Expected host to be not available, got true")
	}
}

func TestIsHostAvailableOnPort_Available(t *testing.T) {
	// Create a listener to simulate a server on localhost:8080
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		t.Error("Failed to create a listener for coverage")
		return
	}
	defer listener.Close()

	// Set the dialerInstance to use the defaultDialer for testing
	dialerInstance = &defaultDialer{}
	defer func() { dialerInstance = &defaultDialer{} }()

	// Use an existing port for testing
	available := isHostAvailableOnPort("localhost", "8080")
	if !available {
		t.Error("Expected host to be available, got false")
	}
}

//func TestShouldRetryRequest_TLSHandshakeFailure_HTTPS(t *testing.T) {
//	err := errors.New("tls: handshake failure")
//	retry := shouldRetryRequest(err, httpsProtocol)
//	if !retry {
//		t.Error("Expected retry for TLS handshake failure, got false")
//	}
//}

//func TestShouldRetryRequest_TLSHandshakeFailure_HTTP(t *testing.T) {
//	err := errors.New("tls: handshake failure")
//	retry := shouldRetryRequest(err, httpProtocol)
//	if retry {
//		t.Error("Expected no retry for TLS handshake failure with HTTP, got true")
//	}
//}

func TestDetermineProtocol_HTTP(t *testing.T) {
	protocol := determineProtocol("example.com:80")
	if protocol != httpProtocol {
		t.Errorf("Expected HTTP protocol, got %s", protocol)
	}
}

func TestDetermineProtocol_HTTPS(t *testing.T) {
	// Assuming the host is not actually available on port 443 in the test environment
	protocol := determineProtocol("example.com")
	if protocol != httpsProtocol {
		t.Errorf("Expected HTTPS protocol, got %s", protocol)
	}
}

func TestCloseConnection(t *testing.T) {
	mockConn := &mockErrorConn{}
	closeConnection(mockConn)
	// Add assertions for the actual closing of the connection
}

func TestDetermineProtocolFromURL_HTTP(t *testing.T) {
	protocol := determineProtocolFromURL("http://example.com")
	if protocol != httpProtocol {
		t.Errorf("Expected HTTP protocol, got %s", protocol)
	}
}

func TestDetermineProtocolFromURL_HTTPS(t *testing.T) {
	protocol := determineProtocolFromURL("https://example.com")
	if protocol != httpsProtocol {
		t.Errorf("Expected HTTPS protocol, got %s", protocol)
	}
}

func TestReplaceProtocol(t *testing.T) {
	newURL := replaceProtocol("http://example.com", httpProtocol, httpsProtocol)
	if newURL != "https://example.com" {
		t.Errorf("Expected URL with HTTPS protocol, got %s", newURL)
	}
}
