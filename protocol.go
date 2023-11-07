package routerosv7_restfull_api

import (
	"log"
	"net"
	"strings"
)

/*
determineProtocol determines the protocol to use for the request.
If the host is available on port 443, HTTPS is used, if not, then HTTP is used.
*/
func determineProtocol(host string) string {
	if isHostAvailableOnPort(host, "443") {
		return httpsProtocol
	}
	return httpProtocol
}

// closeConnection closes a connection.
func closeConnection(conn net.Conn) {
	err := conn.Close()
	if err != nil {
		log.Println(err)
	}
}

// determineProtocolFromURL determines the protocol to use for the request from the URL.
func determineProtocolFromURL(url string) string {
	if strings.HasPrefix(url, httpsProtocol) {
		return httpsProtocol
	}
	return httpProtocol
}

// replaceProtocol replaces the protocol in a URL with a new protocol.
func replaceProtocol(url, oldProtocol, newProtocol string) string {
	return strings.Replace(url, oldProtocol, newProtocol, 1)
}

// ShouldRetryRequest checks if a request should be retried based on the error and protocol.
func shouldRetryRequest(err error, protocol string) bool {
	return strings.Contains(err.Error(), tlsHandshakeFailure) && protocol == httpsProtocol
}

// IsHostAvailableOnPort checks if a host is available on a given port.
func isHostAvailableOnPort(host, port string) bool {
	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		return false
	}
	defer closeConnection(conn)
	return true
}
