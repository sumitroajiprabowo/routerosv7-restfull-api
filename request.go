package routerosv7_restfull_api

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func makeRequest(ctx context.Context, config requestConfig) (interface{}, error) {
	protocol := determineProtocolFromURL(config.URL)
	httpClient := createHTTPClient(protocol)
	requestBody := createRequestBody(config.Payload)
	request, err := createRequest(ctx, config.Method, config.URL, requestBody, config.Username, config.Password)
	if err != nil {
		return nil, err
	}

	response, err := httpClient.Do(request)
	if err != nil {
		if shouldRetryRequest(err, protocol) {
			config.URL = replaceProtocol(config.URL, httpsProtocol, httpProtocol)
			request.URL, _ = request.URL.Parse(config.URL)
			response, err = httpClient.Do(request)
		}
		if err != nil {
			return nil, err
		}
	}

	defer closeResponseBody(response.Body)
	return readJSONResponse(response.Body)
}

func createHTTPClient(protocol string) *http.Client {
	client := &http.Client{}
	if protocol == httpsProtocol {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{},
		}
	}
	return client
}

func createRequestBody(payload []byte) io.Reader {
	if len(payload) > 0 {
		return bytes.NewBuffer(payload)
	}
	return nil
}

func createRequest(ctx context.Context, method, url string, body io.Reader, username, password string) (
	*http.Request, error,
) {
	request, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	request.SetBasicAuth(username, password)
	return request, nil
}

func closeResponseBody(body io.ReadCloser) {
	err := body.Close()
	if err != nil {
		log.Println(err)
	}
}

func readJSONResponse(body io.ReadCloser) (interface{}, error) {
	var responseData interface{}
	err := json.NewDecoder(body).Decode(&responseData)
	return responseData, err
}
