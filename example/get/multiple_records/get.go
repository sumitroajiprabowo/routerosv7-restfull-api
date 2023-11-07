package main

import (
	"context"
	"encoding/json"
	"fmt"
	pkg "github.com/megadata-dev/routerosv7-restfull-api"
)

const (
	defaultHost     = "10.90.0.251"
	defaultUsername = "userapi"
	defaultPassword = "password"
	defaultCommand  = "ip/address"
)

type AppConfig struct {
	Host     string
	Username string
	Password string
	Command  string
}

func NewAppConfig(host, username, password, command string) *AppConfig {
	return &AppConfig{
		Host:     host,
		Username: username,
		Password: password,
		Command:  command,
	}
}

type DataRetriever interface {
	GetData(ctx context.Context) (interface{}, error)
}

type RouterOSDataRetriever struct {
	config *AppConfig
}

func NewRouterOSDataRetriever(config *AppConfig) DataRetriever {
	return &RouterOSDataRetriever{config: config}
}

func (r *RouterOSDataRetriever) GetData(ctx context.Context) (interface{}, error) {
	data, err := pkg.Print(ctx, r.config.Host, r.config.Username, r.config.Password, r.config.Command)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func PrintJSON(data interface{}) {
	response, ok := data.([]interface{})
	if !ok {
		fmt.Println("Failed to process data")
		return
	}

	// Process and print the data as JSON
	var ipAddresses []map[string]interface{}
	for _, item := range response {
		if ipItem, ok := item.(map[string]interface{}); ok {
			ipAddresses = append(ipAddresses, ipItem)
		}
	}

	jsonData, err := json.Marshal(map[string]interface{}{
		"code":   200,
		"status": "OK",
		"data":   ipAddresses,
	})
	if err != nil {
		fmt.Println("Failed to marshal JSON:", err)
		return
	}

	fmt.Println(string(jsonData))
}

func main() {
	config := NewAppConfig(defaultHost, defaultUsername, defaultPassword, defaultCommand)
	pingManager := pkg.NewPing(config.Host)

	if pingManager == nil {
		fmt.Println("Failed to create PingManager")
		return
	}

	err := pingManager.CheckAvailableDevice()
	if err != nil {
		fmt.Println("Device is not available:", err)
	} else {
		fmt.Println("Device is available")
	}

	err = pkg.AuthDevice(context.Background(), config.Host, config.Username, config.Password)
	if err != nil {
		fmt.Println("Authentication failed:", err)
	} else {
		fmt.Println("Authentication success")
	}

	dataRetriever := NewRouterOSDataRetriever(config)
	data, err := dataRetriever.GetData(context.Background())
	if err != nil {
		fmt.Println("Failed to retrieve data:", err)
		return
	}

	PrintJSON(data)
}
