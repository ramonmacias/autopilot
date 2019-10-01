package apiAutopilot

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"time"
)

var (
	c    *configuration
	once sync.Once
)

const (
	customAutopilotAuthorizationHeader = "autopilotapikey"
)

type APIClientable interface {
	GetClient() *http.Client
	GetBaseUrl() string
	SendRequest(method, url, authToken string, body io.Reader) (*ContactResponse, error)
}

type configuration struct {
	client  http.Client
	baseUrl string
}

type configInfo struct {
	BaseUrl string `json:"base_url"`
	TimeOut int    `json:"time_out_seconds"`
}

func NewApiClient() *configuration {
	return autopilotConfig()
}

func autopilotConfig() *configuration {
	once.Do(func() {
		c = setupConfig()
	})
	return c
}

func setupConfig() *configuration {
	path, err := filepath.Abs("../../config/autopilot_client.json")
	if err != nil {
		log.Printf("Error while try to get abs path: %v", err)
	}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panicf("There is an error while try to read autopilot config file: %v", err)
	}
	configInfo := &configInfo{}
	if err = json.Unmarshal([]byte(file), configInfo); err != nil {
		log.Panicf("There is an error while try to unmarshal the json autopilot client config info, err: %v", err)
	}
	timeout := time.Duration(time.Duration(configInfo.TimeOut) * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	return &configuration{
		client:  client,
		baseUrl: configInfo.BaseUrl,
	}
}

func (c *configuration) GetClient() *http.Client {
	return &c.client
}

func (c *configuration) GetBaseUrl() string {
	return c.baseUrl
}

type ContactResponse struct {
	Email string `json:"Email"`
	Id    string `json:"contact_id"`
	Body  []byte `json:"-"`
}

type Error struct {
	StatusCode int
	Message    string
}

func (e *Error) Error() string {
	return fmt.Sprintf("Error from autopìlot API with code: %d and body: %s", e.StatusCode, e.Message)
}

func (c *configuration) SendRequest(method, url, authToken string, body io.Reader) (*ContactResponse, error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Printf("There is an error while try to create a request for create contact, err: %v", err)
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set(customAutopilotAuthorizationHeader, authToken)

	resp, err := c.client.Do(request)
	responseBody, _ := ioutil.ReadAll(resp.Body)
	if err != nil || resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &Error{
			StatusCode: resp.StatusCode,
			Message:    string(responseBody),
		}
	}

	response := &ContactResponse{}
	json.Unmarshal(responseBody, response)
	response.Body = responseBody
	return response, nil
}
