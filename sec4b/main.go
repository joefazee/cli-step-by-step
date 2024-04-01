package main

import (
	"fmt"
	"net/http"
	"time"
)

type SiteConfig struct {
	URL             string
	AcceptableCodes []int
	Frequency       int
}

type Result struct {
	URL    string
	Up     bool
	Status int
}

type HttpClient interface {
	Get(url string) (resp *http.Response, err error)
}

type DefaultHttpClient struct{}

func (c *DefaultHttpClient) Get(url string) (resp *http.Response, err error) {
	return http.Get(url)
}

func check(config SiteConfig, client HttpClient, results chan<- Result) {
	resp, err := client.Get(config.URL)

	result := Result{
		URL: config.URL,
	}

	if err != nil {
		result.Up = false
		results <- result
		return
	}

	defer resp.Body.Close()
	result.Status = resp.StatusCode
	result.Up = false

	for _, code := range config.AcceptableCodes {
		if resp.StatusCode == code {
			result.Up = true
			break
		}
	}

	results <- result
}

func scheduleCheck(config SiteConfig, client HttpClient, results chan<- Result) {
	go func() {
		ticker := time.NewTicker(time.Duration(config.Frequency) * time.Second)
		for {
			select {
			case <-ticker.C:
				check(config, client, results)
			}
		}
	}()
}

func main() {

	sites := []SiteConfig{
		{
			URL:             "https://google.com",
			AcceptableCodes: []int{200},
			Frequency:       5,
		},

		{
			URL:             "https://yahoo.com",
			AcceptableCodes: []int{200},
			Frequency:       6,
		},

		{
			URL:             "https://go.dev",
			AcceptableCodes: []int{200},
			Frequency:       7,
		},

		{
			URL:             "https://localhost.com.me",
			AcceptableCodes: []int{200},
			Frequency:       2,
		},
	}

	results := make(chan Result)
	client := &DefaultHttpClient{}

	for _, site := range sites {
		scheduleCheck(site, client, results)
	}

	for result := range results {
		if result.Up {
			fmt.Printf("%s is up (Status code: %d)\n", result.URL, result.Status)
		} else {
			fmt.Printf("%s is down (Status code: %d)\n", result.URL, result.Status)
		}
	}
}
