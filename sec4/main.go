package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
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

type SiteMonitor struct {
	configs    []SiteConfig
	results    chan Result
	client     HttpClient
	configFile string
	lock       sync.RWMutex
}

type yamlSiteConfig struct {
	URL             string `yaml:"url"`
	AcceptableCodes []int  `yaml:"acceptableCodes"`
	Frequency       int    `yaml:"frequency"`
}

func (m *SiteMonitor) loadConfigs() error {
	file, err := os.ReadFile(m.configFile)
	if err != nil {
		return err
	}

	var yamlConfigs []yamlSiteConfig
	err = yaml.Unmarshal(file, &yamlConfigs)
	if err != nil {
		return err
	}

	m.configs = make([]SiteConfig, len(yamlConfigs))
	for i, c := range yamlConfigs {
		m.configs[i] = SiteConfig{
			URL:             c.URL,
			AcceptableCodes: c.AcceptableCodes,
			Frequency:       c.Frequency,
		}
	}

	return nil
}

func (m *SiteMonitor) saveConfigs() error {
	m.lock.RLock()
	defer m.lock.RUnlock()

	yamlConfigs := make([]yamlSiteConfig, len(m.configs))
	for i, config := range m.configs {
		yamlConfigs[i] = yamlSiteConfig{
			URL:             config.URL,
			AcceptableCodes: config.AcceptableCodes,
			Frequency:       config.Frequency,
		}
	}

	data, err := yaml.Marshal(yamlConfigs)
	if err != nil {
		return err
	}

	return os.WriteFile(m.configFile, data, 0466)

}
func main() {
	configFile := "./sec4/sites.yaml"

	monitor := &SiteMonitor{
		configFile: configFile,
		results:    make(chan Result),
		client:     &DefaultHttpClient{},
	}

	err := monitor.loadConfigs()
	if err != nil {
		log.Fatal(err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sig
		fmt.Printf("we got the signal %s\n", sig)
		err := monitor.saveConfigs()
		if err != nil {
			log.Println("Error", err.Error())
		}
		os.Exit(0)
	}()

	for _, site := range monitor.configs {
		scheduleCheck(site, monitor.client, monitor.results)
	}

	for result := range monitor.results {
		if result.Up {
			fmt.Printf("%s is up (Status code: %d)\n", result.URL, result.Status)
		} else {
			fmt.Printf("%s is down (Status code: %d)\n", result.URL, result.Status)
		}
	}
}
