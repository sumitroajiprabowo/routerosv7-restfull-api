package routerosv7_restfull_api

import (
	"log"
	"net"
)

func determineProtocol(host string) (string, error) {
	conn, err := net.Dial("tcp", host+":443")
	if err != nil {
		return "http", nil
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println(err)
		}
	}(conn)
	return "https", nil
}
