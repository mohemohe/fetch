package fetch

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
)

type (
	testFetchResult struct {
		Args struct {
			Foo1 string `json:"foo1"`
			Foo2 string `json:"foo2"`
		} `json:"args"`
		Headers map[string]string `json:"headers"`
	}

	testFetchPostBody struct {
		Key string `json:"key"`
	}
	testFetchPostResult struct {
		Json testFetchPostBody `json:"json"`
	}
)

func TestFetch(t *testing.T) {
	url := "https://postman-echo.com/get?foo1=bar1&foo2=bar2"
	headerKey := "x-test-header"
	headerValue := "This is test text."
	option := Option{
		Headers: map[string]string{
			headerKey: headerValue,
		},
	}
	req, err := Fetch(url, option)
	if err != nil {
		t.Error(err)
	}
	if req.Response() == nil {
		t.Error(errors.New("invalid state"))
	}
	if !req.OK {
		t.Error(errors.New("remote server response error"))
	}
	if req.Status != 200 {
		t.Error(errors.New("invalid status code"), req.Status)
	}
	if req.StatusText != "200" {
		t.Error(errors.New("invalid status text"), req.StatusText)
	}

	result := new(testFetchResult)
	if err := req.JSON(result); err != nil {
		t.Error(err)
	}
	if result.Headers[headerKey] != headerValue {
		t.Error(errors.New("invalid response header"), *result)
	}
	if result.Args.Foo1 != "bar1" || result.Args.Foo2 != "bar2" {
		t.Error(errors.New("invalid response body"), *result)
	}
}

func TestFetchMap(t *testing.T) {
	url := "https://postman-echo.com/get?foo1=bar1&foo2=bar2"
	req, err := Fetch(url)
	if err != nil {
		t.Error(err)
	}

	result, err := req.Map()
	if err != nil {
		t.Error(err)
	}
	if result["args"].(map[string]interface{})["foo1"] != "bar1" || result["args"].(map[string]interface{})["foo2"] != "bar2" {
		t.Error(errors.New("invalid response"), result)
	}
}

func TestFetchPost(t *testing.T) {
	value := "This is test text"
	body := testFetchPostBody{
		Key: value,
	}
	b, err := json.Marshal(body)
	if err != nil {
		t.Error(err)
	}
	option := Option{
		Method: http.MethodPost,
		Body:   bytes.NewReader(b),
		Headers: map[string]string{
			"content-type": "application/json",
		},
	}
	url := "https://postman-echo.com/post"
	req, err := Fetch(url, option)
	if err != nil {
		t.Error(err)
	}
	if req.Response() == nil {
		t.Error(errors.New("invalid state"))
	}

	result := new(testFetchPostResult)
	if err := req.JSON(result); err != nil {
		t.Error(err)
	}
	if result.Json.Key != value {
		t.Error(errors.New("invalid response"), result)
	}
}
