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

// main function for this example application to delete data from RouterOS device
func main() {

	// ipAddress variable for the IP address to be deleted
	ipAddress := "192.168.99.1/24"

	// Check if the address already exists in the RouterOS device
	exists, err := checkData(routerIP, username, password, "ip/address", "address", ipAddress)
	if err != nil {
		fmt.Println("Failed to check data:", err)
		return
	}

	// Print error message if the address does not exist and return from this function
	if !exists {
		jsonErrorNotFound := webResponse{
			Code:   404,
			Status: "Not Found",
			Data:   "Address not found",
		}

		// Marshal the jsonErrorNotFound to JSON
		jsonData, err := json.Marshal(jsonErrorNotFound)

		// Print error message if there is an error and return from this function
		if err != nil {
			fmt.Println("Failed to marshal JSON:", err)
			return
		}

		// Print the JSON string to the console
		fmt.Println(string(jsonData))

		return
	}

	// Perform the DELETE operation
	command := fmt.Sprintf("ip/address/%s", getAddressID(routerIP, username, password, "ip/address", "address", ipAddress))

	// Delete the data
	_, err = deleteData(routerIP, username, password, command)
	if err != nil {
		fmt.Println("Failed to delete data:", err)
		return
	}

	// Create jsonSuccess variable as webResponse struct
	jsonSuccess := webResponse{
		Code:   204,
		Status: "No Content",
		Data:   nil,
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

}

// checkData function to check if the data exists in the RouterOS device
func checkData(routerIP, username, password, command, field, value string) (bool, error) {

	// Retrieve the data from the RouterOS device using the Print function
	data, err := routerosv7_restfull_api.Print(context.Background(), routerIP, username, password, command)
	if err != nil {
		return false, err
	}

	// Check if the field with the given value exists in the retrieved data
	response, ok := data.([]interface{})
	if ok {
		for _, item := range response {
			if dataItem, ok := item.(map[string]interface{}); ok {
				if fieldValue, ok := dataItem[field].(string); ok {
					if fieldValue == value {
						return true, nil
					}
				}
			}
		}
	}

	return false, nil
}

// getAddressID function to get the ID of the address to be deleted
func getAddressID(routerIP, username, password, command, field, value string) string {

	// Retrieve the data from the RouterOS device using the Print function
	data, err := routerosv7_restfull_api.Print(context.Background(), routerIP, username, password, command)
	if err != nil {
		return ""
	}

	// Check if the field with the given value exists in the retrieved data
	response, ok := data.([]interface{})
	if ok {
		for _, item := range response {
			if dataItem, ok := item.(map[string]interface{}); ok {
				if fieldValue, ok := dataItem[field].(string); ok {
					if fieldValue == value {
						if id, ok := dataItem[".id"].(string); ok {
							return id
						}
					}
				}
			}
		}
	}

	return ""
}

// deleteData function to delete data from RouterOS device
func deleteData(routerIP, username, password, command string) (interface{}, error) {
	data, err := routerosv7_restfull_api.Delete(context.Background(), routerIP, username, password, command)
	if err != nil {
		return err, nil
	}

	return data, nil
}
