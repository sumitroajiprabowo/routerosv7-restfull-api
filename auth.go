package routerosv7_restfull_api

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func AuthDevice(ctx context.Context, host, username, password string) error {
	// Determine the protocol from the URL (HTTP or HTTPS) and check if the protocol is valid
	// If the protocol is not valid, return error
	protocol := determineProtocol(host)

	// Create the URL for the request to Mikrotik Router
	url := fmt.Sprintf("%s://%s/rest/system/resource", protocol, host)

	// Create request to Mikrotik Router based on "url" and "method" (GET)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	// Set basic auth for GET request to Mikrotik Router
	req.SetBasicAuth(username, password)

	// Create an HTTP client
	var client *http.Client

	// Check if protocol is https or http
	if protocol == httpsProtocol {
		// Attempt HTTPS, but if there is a TLS handshake error, fall back to HTTP
		client = &http.Client{}
	} else {
		// Use a client with insecure TLS configuration for HTTP
		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
	}

	req = req.WithContext(ctx)

	// Send the request
	resp, err := client.Do(req)

	if err != nil {
		if strings.Contains(err.Error(), "tls: handshake failure") {
			protocolWeb := httpProtocol
			// TLS handshake failure, redirect to HTTP
			url = fmt.Sprintf("%s://%s/rest/system/resource", protocolWeb, host)
			req.URL, _ = req.URL.Parse(url)
			resp, err = client.Do(req) // Retry with HTTP
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	// Close response body
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)

	// Check if the response status code is 404 or 401 return exception.Unauthorized
	if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusUnauthorized {
		return errors.New("unauthorized")
	}

	return nil
}
