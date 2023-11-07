package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/megadata-dev/routerosv7-restfull-api"
)

// Create constants for the default values for this example application
const (
	routerIP = "10.90.0.251" // Change this to your router's IP address
	username = "userapi"     // Change this to your router's username
	password = "password"    // Change this to your router's password
)

// webResponse struct for web response data
type webResponse struct {
	Code   int         `json:"code"`   // HTTP status code
	Status string      `json:"status"` // HTTP status text
	Data   interface{} `json:"data"`   // Data to be sent
}

// checkIfAddressExists checks if the address already exists
func checkIfAddressExists(routerIP, username, password, address string) (bool, error) {

	// Get Data with Print function
	data, err := routerosv7_restfull_api.Print(context.Background(), routerIP, username, password, "ip/address")

	// Check if there is an error
	if err != nil {
		return false, err
	}

	// Check if the data is an array of map[string]interface{}
	response, ok := data.([]interface{})
	if ok {
		for _, item := range response {
			if dataItem, ok := item.(map[string]interface{}); ok {
				if fieldValue, ok := dataItem["address"].(string); ok {
					if fieldValue == address {
						return true, nil
					}
				}
			}
		}
	}

	// Return false if the address does not exist
	return false, nil
}

// addAddress adds an address to the router and returns the response
func addAddress(routerIP, username, password, command string, payload []byte) (map[string]interface{}, error) {

	// Add data with AddData function
	response, err := routerosv7_restfull_api.AddData(context.Background(), routerIP, username, password, command, payload)

	// Check if there is an error
	if err != nil {
		return nil, err
	}

	// Return the response
	return response.(map[string]interface{}), nil
}

// getAddressByID gets an address by ID for checking the newly added address data already exists or not
func getAddressByID(routerIP, username, password, id string) (interface{}, error) {

	// Get data with Print function
	data, err := routerosv7_restfull_api.Print(context.Background(), routerIP, username, password, "ip/address/"+id)

	// Check if there is an error
	if err != nil {
		return nil, err
	}

	// Return the data as interface{} and nil error
	return data, nil
}

// main function for this example application
func main() {

	// Create params variable as map[string]string
	ipAddress := "192.168.99.1/24" // Change this to your desired IP address
	iface := "ether1"              // Change this to your desired interface

	// Check if the address already exists
	if exists, err := checkIfAddressExists(routerIP, username, password, ipAddress); err != nil {
		fmt.Println("Failed to check if address exists:", err)
		return
	} else if exists {
		// Create jsonError variable
		jsonError := webResponse{
			Code:   409,
			Status: "Conflict",
			Data:   fmt.Sprintf("Address %s already exists", ipAddress),
		}

		// Marshal the jsonError to JSON
		jsonData, err := json.Marshal(jsonError)

		// Print error message if there is an error and return from this function
		if err != nil {
			fmt.Println("Failed to marshal JSON:", err)
			return
		}

		// Print the JSON string to the console
		fmt.Println(string(jsonData))
		return
	}

	// Create payload variable as []byte with the desired payload data
	payload := fmt.Sprintf(`{"address": "%s","interface": "%s"}`, ipAddress, iface)

	// Add address with addAddress function and get the response data as map[string]interface{} if there is no error
	//and print the response data to the console as JSON string
	if response, err := addAddress(routerIP, username, password, "ip/address", []byte(payload)); err != nil {
		fmt.Println("Failed to add address:", err)
		return
	} else if id, ok := response["ret"].(string); ok {
		data, err := getAddressByID(routerIP, username, password, id)
		if err != nil {
			fmt.Println("Failed to get new address added:", err)
			return
		}

		// Create jsonSuccess variable
		jsonSuccess := webResponse{
			Code:   200,
			Status: "OK",
			Data:   data,
		}

		// Marshal the jsonSuccess to JSON
		jsonData, err := json.Marshal(jsonSuccess)

		// Print error message if there is an error and return from this function
		if err != nil {
			fmt.Println("Failed to marshal JSON:", err)
			return
		}

		// Print the JSON string to the console
		fmt.Println(string(jsonData))
	} else {
		fmt.Println("Failed to get new address added")
	}
}
