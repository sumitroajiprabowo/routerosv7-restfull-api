## RouterOS v7 Rest API Client 
## Description
This library is a wrapper for the RouterOS v7 Rest API. It provides a simple way to access RouterOS devices using Go.

## HTTP Methods
| HTTP| RouterOS | Description                                             |
| --- |----------|---------------------------------------------------------|
| GET | print    | Get data record                                         |
| PUT | add      | Create a new record                                     |
|PATCH| set      | Update a record                                         |
|DELETE| remove  | Delete a record                                         |
| POST|      | Universal method to get access to all console commands  |

for more info see: https://help.mikrotik.com/docs/display/ROS/REST+API#RESTAPI-Overview

## Overview
Package routerosv7_restfull_api functions:
- **NewPing** - function to create new PingManager instance
- **CheckAvailableDevice** - function to check if the device is available
- **Auth** - function to authenticate the Mikrotik device
- **Print** - function to get data record
- **Add** - function to create a new record
- **Set** - function to update a record
- **Remove** - function to delete a record
- **Run** - function to run console commands
- **Exec** - function to execute the request to Mikrotik device

## Usage
### Ping Device
This is example implementation check if the device is available and get data record
```go
pingManager := routerosv7_restfull_api.NewPing(host)

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
````
### Auth
This is example implementation to authenticate the Mikrotik device
```go
request, err := routerosv7_restfull_api.Auth(context.Background(), 
	routerosv7_restfull_api.AuthDeviceConfig{
        Host:     "192.168.88.1",
        Username: "admin",
        Password: "changeme",
    }
)
````

### Print
This is example implementation to get data record
```go
request := routerosv7_restfull_api.Print(
    "192.168.88.1",
    "admin",
    "changeme",
    "ip/address"
)
```

### Add
This is example implementation to create a new record
```go
request := routerosv7_restfull_api.Add(
    "192.168.88.1",
    "admin",
    "changeme",
    "ip/address",
    []byte(`{
        "address": "192.168.100.1/24",
        "interface": "ether1"}`),
)
```

### Set
This is example implementation to update a record
```go
request := routerosv7_restfull_api.Set(
    "192.168.88.1",
    "admin",
    "changeme",
    "ip/address/*5F",
    []byte(`{"comment": "Test API"}`),
)
```

### Remove
This is example implementation to delete a record
```go
request := routerosv7_restfull_api.Remove(
    "192.168.88.1",
    "admin",
    "changeme",
    "ip/address/*5F",
)
```

### Run
This is example implementation to run console commands
###### Print
```go
request := routerosv7_restfull_api.Run(
		"192.168.88.1",
		"admin",
		"changeme",
		"ip/address/print",
		nil
)
```

#### Add
```go
request := routerosv7_restfull_api.Run(
		"192.168.88.1",
		"admin",
		"changeme",
		"ip/address/add",
		[]byte(`{
        "address": "192.168.100.1/24",
        "interface": "ether1"}`),
)
```

#### Proplist
```go
request := routerosv7_restfull_api.Run(
    "192.168.88.1",
    "admin",
    "changeme",
    "ip/address/print",
    []byte(`{"_proplist": ["address","interface"]}`),
)
```

#### Query
```go
request := routerosv7_restfull_api.Run(
		"192.168.88.1",
		"admin",
		"changeme",
		"ip/address/print",
		[]byte(`{
			"_proplist": ["address","interface"],
			".query": ["network=192.168.100.0","#"]
		}`),
)
```

### SSL Certificates on RouterOS v7
You have to set up certificates to use secure HTTPS
More info here: https://help.mikrotik.com/docs/display/ROS/Certificates
```shell
/certificate
add name=ca-template days-valid=3650 common-name=your.server.url key-usage=key-cert-sign,crl-sign
add name=server-template days-valid=3650 common-name=your.server.url

/certificate
sign ca-template name=root-ca
:delay 3s
sign ca=root-ca server-template name=server
:delay 3s

/certificate
set root-ca trusted=yes
set server trusted=yes

/ip service
set www-ssl certificate=server disabled=no
```