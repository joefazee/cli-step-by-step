package main

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockHttpClient struct {
	Response *http.Response
	Error    error
}

func (c *MockHttpClient) Get(url string) (resp *http.Response, err error) {
	return c.Response, c.Error
}

func Test_check(t *testing.T) {
	mockClient := &MockHttpClient{
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString("ok")),
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

	result := <-results

	if !result.Up || result.Status != 200 {
		t.Error("site should be up with 200 status code")
	}
}

func Test_check_Error(t *testing.T) {
	mockClient := &MockHttpClient{
		Error: errors.New("error"),
	}

	config := SiteConfig{
		URL: "http://localhost.com",
	}

	results := make(chan Result, 1)

	check(config, mockClient, results)

	result := <-results

	if result.Up {
		t.Error("site should be down")
	}
}

func TestDefaultClient_Get(t *testing.T) {

	client := DefaultClient{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	}))
	defer srv.Close()

	resp, err := client.Get(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Error("expect status code to be 200")
	}

	expectedBody := "ok"
	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if expectedBody != string(bs) {
		t.Errorf("expect %v; got %v", expectedBody, string(bs))
	}

}
