Utilite for logging and add default custom http headers into http client

This simple example shows how you can add default header for http client and setup log into default logger.

```go
package main

import (
    "github.com/ciricc/httpaddtransport"
    "http"
)

func main() {
    httpAddTransport, err := httpaddtransport.New(nil, nil)
    
    if err != nil {
		panic(err)
	}

    httpAddTransport.Headers.Set("User-Agent", "go/1.19")
    httpAddTransport.Log = true

    httpClient := http.Client{
        Transport: httpAddTransport
    }

	_, err = httpClient.Get("https://google.com/")
	if err != nil {
		panic(err)
	}
}

```