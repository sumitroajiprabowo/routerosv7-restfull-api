package routerosv7_restfull_api

import (
	"log"
	"net"
)

func determineProtocol(host string) (string, error) {
	conn, err := net.Dial("tcp", host+":443")
	if err != nil {
		return httpProtocol, nil
	}
	defer closeConnection(conn)
	return httpsProtocol, nil
}

func closeConnection(conn net.Conn) {
	err := conn.Close()
	if err != nil {
		log.Println(err)
	}
}
