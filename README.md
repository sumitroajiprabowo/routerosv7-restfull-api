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
- **Auth** - function to authenticate the Mikrotik device
- **Print** - function to get data record
- **Add** - function to create a new record
- **Set** - function to update a record
- **Remove** - function to delete a record
- **Run** - function to run console commands

## Usage
### Auth
This is example implementation to authenticate the Mikrotik device
```go
package main

import (
	"context"
	"fmt"
	"github.com/sumitroajiprabowo/routerosv7-restfull-api"
)

func main() {
    // Authenticate to the router
    _, err := routerosv7_restfull_api.Auth(context.Background(),
        routerosv7_restfull_api.AuthConfig{
            Host:     "192.168.88.1", // Change this to your router's IP address
            Username: "username",      // Change this to your router's username
            Password: "password",     // Change this to your router's password
        })
    // Check if authentication failed
    if err != nil {
        fmt.Println(err)
    } else {
        fmt.Println("Authentication success")
    }
}
````

### Print
This is example implementation to get data record
```go
package main

import (
	"context"
	"fmt"
	"github.com/sumitroajiprabowo/routerosv7-restfull-api"
)

func main() {
	data, err := routerosv7_restfull_api.Print(
		context.Background(),
		"192.168.88.1",                     // Change this to your router's IP address
		"username",                         // Change this to your router's username
		"password",                         // Change this to your router's password
		"ip/address",                       // Change this to the path you want to add
	)
	
	if err != nil {
		fmt.Println("Failed to print:", err)
		return
	}
	
	fmt.Println(data)
	return
}
```

### Add
This is example implementation to create a new record
```go
package main

import (
	"context"
	"fmt"
	"github.com/sumitroajiprabowo/routerosv7-restfull-api"
)

func main() {
	payload := fmt.Sprintf(`{"address": "%s","interface": "%s"}`, "192.168.99.1/24", "ether1")
	data, err := routerosv7_restfull_api.Add(
		context.Background(),
		"192.168.88.1",                     // Change this to your router's IP address
		"username",                         // Change this to your router's username
		"password",                         // Change this to your router's password
		"ip/address",                       // Change this to the path you want to add
		[]byte(payload))                    // Payload to add

	if err != nil {
		fmt.Println("Failed to add address:", err)
		return
	}

	fmt.Println(data)
	return
}
```

### Set
This is example implementation to update a record
```go
package main

import (
	"context"
	"fmt"
	"github.com/sumitroajiprabowo/routerosv7-restfull-api"
)

func main() {
	data, err := routerosv7_restfull_api.Set(
		context.Background(),
		"192.168.88.1",                    // Change this to your router's IP address
		"username",                        // Change this to your router's username
		"password",                        // Change this to your router's password
		"ip/address/*7A",                  // *7A is the ID of the address
		[]byte(`{"comment": "Test API"}`), // Change this to the payload you want to set
	)
	
	if err != nil {
		fmt.Println("Failed to patch data:", err)
		return
	}
	
	fmt.Println(data)
	return
}
```

### Remove
This is example implementation to delete a record
```go
package main

import (
	"context"
	"fmt"
	"github.com/sumitroajiprabowo/routerosv7-restfull-api"
)

func main() {
	if data, _ := routerosv7_restfull_api.Remove(
		context.Background(),
		"192.168.88.1",   // Change this to your router's IP address
		"username",       // Change this to your router's username
		"password",       // Change this to your router's password
		"ip/address/*7A", // *7A is the ID of the address
	); data != nil {
		fmt.Println("Failed to delete data:", data)
		return
	} else {
		fmt.Println("Data successfully deleted")
	}
}
```

### Run
This is example implementation to run console commands
###### Print
```go
package main

import (
	"context"
	"fmt"
	"github.com/sumitroajiprabowo/routerosv7-restfull-api"
)

func main() {
	data, err := routerosv7_restfull_api.Run(
		context.Background(), // Create a context
		"192.168.88.1",       // Change this to your router's IP address
		"username",           // Change this to your router's username
		"password",           // Change this to your router's password
		"ip/address/print",   // Change this to the command you want to execute
		nil,                  
	)
	
	if err != nil {
		fmt.Println("Failed to retrieve data:", err)
		return
	}
	
	fmt.Println(data)
}
```

#### Add
```go
package main

import (
	"context"
	"fmt"
	"github.com/sumitroajiprabowo/routerosv7-restfull-api"
)

func main() {
	data, err := routerosv7_restfull_api.Run(
		context.Background(), // Create a context
		"192.168.88.1",       // Change this to your router's IP address
		"username",            // Change this to your router's username
		"password",           // Change this to your router's password
		"ip/address/add",     // Change this to the command you want to execute
		[]byte(`{
        "address": "192.168.100.1/24",
        "interface": "ether1"}`), // Change this to the payload you want to send
	)
	
	if err != nil {
		fmt.Println("Failed to add address:", err)
		return
	}
	
	fmt.Println(data)
}
```

#### Proplist
```go
package main

import (
	"context"
	"fmt"
	"github.com/sumitroajiprabowo/routerosv7-restfull-api"
)

func main() {
	data, err := routerosv7_restfull_api.Run(
		context.Background(), // Create a context
		"192.168.88.1",       // Change this to your router's IP address
		"username",            // Change this to your router's username
		"password",           // Change this to your router's password
		"ip/address/print",   // Change this to the command you want to execute
		[]byte(`{"_proplist": ["address","interface"]}`), // Change this to the desired payload data
	)
	
	if err != nil {
		fmt.Println("Failed to execute command:", err)
		return
	}
	
	fmt.Println(data)
}
```

#### Query
```go
package main

import (
	"context"
	"fmt"
	"github.com/sumitroajiprabowo/routerosv7-restfull-api"
)

func main() {
	data, err := routerosv7_restfull_api.Run(
		context.Background(), // Create a context
		"192.168.88.1",       // Change this to your router's IP address
		"username",           // Change this to your router's username
		"password",           // Change this to your router's password
		"ip/address/print",   // Change this to the command you want to execute
		[]byte(`{
			"_proplist": ["address","interface"],
			".query": ["network=192.168.100.0","#"]
		}`)) // Change this to the payload you want to send
	
	if err != nil {
		fmt.Println("Failed to execute command:", err)
		return
	}
	
	fmt.Println(data)
}
```

### SSL Certificates on RouterOS v7
To use secure HTTPS, set up certificates on RouterOS v7. More information can be found [here](https://help.mikrotik.com/docs/display/ROS/Certificates).

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
