package api

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var (
	client http.Client
)

const (
	customAutopilotAuthorizationHeader = "autopilotapikey"
	autopilotBaseContactURL            = "https://api2.autopilothq.com/v1/contact"
)

func init() {
	timeout := time.Duration(5 * time.Second)
	client = http.Client{
		Timeout: timeout,
	}
}

func ShowContact(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func UpdateContact(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
}

func CreateContact(w http.ResponseWriter, r *http.Request) {

	request, err := http.NewRequest("POST", autopilotBaseContactURL, r.Body)
	if err != nil {
		log.Printf("There is an error while try to create a request for create contact, err: %v", err)
		w.WriteHeader(http.StatusBadGateway)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set(customAutopilotAuthorizationHeader, r.Header.Get(customAutopilotAuthorizationHeader))

	resp, err := client.Do(request)
	if err != nil {
		log.Printf("There is an error while try to send a create contact request to Autopilot API, err : %v", err)
		w.WriteHeader(http.StatusBadGateway)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Unexpected error: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
