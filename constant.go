package routerosv7_restfull_api

import (
	"time"
)

const (
	httpProtocol        = "http"
	httpsProtocol       = "https"
	tlsHandshakeFailure = "tls: handshake failure"
	pingCount           = 3
	pingTimeout         = 1000 * time.Millisecond
	pingInterval        = 100 * time.Millisecond
)

type requestConfig struct {
	URL      string
	Method   string
	Payload  []byte
	Username string
	Password string
}
