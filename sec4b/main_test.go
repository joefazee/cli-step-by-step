package main

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"
)

type MockHttpClient struct {
	Response *http.Response
	Error    error
}

func (c *MockHttpClient) Get(url string) (resp *http.Response, err error) {
	return c.Response, c.Error
}

func Test_Check(t *testing.T) {

	mockClient := &MockHttpClient{
		Response: &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString("OK")),
		},
		Error: nil,
	}

	config := SiteConfig{
		URL:             "http://localhost.com",
		AcceptableCodes: []int{200},
		Frequency:       1,
	}

	results := make(chan Result, 1)

	check(config, mockClient, results)

	res := <-results

	if !res.Up || res.Status != 200 {
		t.Errorf("expect site to be up with 200 status code")
	}
}

func Test_Check_Error(t *testing.T) {

	mockClient := &MockHttpClient{
		Error: errors.New("error"),
	}

	config := SiteConfig{
		URL:             "http://localhost.com",
		AcceptableCodes: []int{200},
		Frequency:       1,
	}

	results := make(chan Result, 1)

	check(config, mockClient, results)

	res := <-results

	if res.Up {
		t.Errorf("expect site down but site is up")
	}
}
