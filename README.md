# fetch-go

JavaScript like http request wrapper for golang

## usage

```sh
go get -u github.com/mohemohe/fetch
```

```go
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type (
	testFetchResult struct {
		Args struct {
			Foo1 string `json:"foo1"`
			Foo2 string `json:"foo2"`
		} `json:"args"`
	}
)

func main() {
	url := "https://postman-echo.com/get?foo1=bar1&foo2=bar2"
	req, err := Fetch(url)
	if err != nil {
		panic(err)
	}
	if !req.OK {
		panic("remote server error")
	}

	result := new(testFetchResult)
	if err := req.JSON(result); err != nil {
		panic(err)
	}
    
    println(result.Args.Foo1)
    // => bar1
}
```