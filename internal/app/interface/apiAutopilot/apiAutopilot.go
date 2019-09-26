package apiAutopilot

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ramonmacias/autopilot/internal/app/domain/model"
)

const (
	customAutopilotAuthorizationHeader = "autopilotapikey"
	autopilotBaseContactURL            = "https://api2.autopilothq.com/v1/contact"
)

type ContactResponse struct {
	Email string `json:"Email"`
	Id    string `json:"contact_id"`
}

type apiAutopilot struct {
	client http.Client
}

func NewApiAutopilot(client http.Client) *apiAutopilot {
	return &apiAutopilot{
		client: client,
	}
}

func (a *apiAutopilot) GetContact(id, authToken string) (*model.Contact, error) {
	request, err := http.NewRequest("GET", autopilotBaseContactURL+"/"+id, nil)
	if err != nil {
		log.Printf("There is an error while try to create a request for create contact, err: %v", err)
		return nil, err
	}
	request.Header.Set(customAutopilotAuthorizationHeader, authToken)

	resp, err := a.client.Do(request)
	if err != nil {
		log.Printf("There is an error while try to send a create contact request to Autopilot API, err : %v", err)
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Unexpected error: %v", err)
		return nil, err
	}
	resp.Body.Close()
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	contactResponse := &ContactResponse{}
	json.NewDecoder(resp.Body).Decode(contactResponse)

	return &model.Contact{
		Id:    contactResponse.Id,
		Email: contactResponse.Email,
		Data:  string(body),
	}, nil
}

func (a *apiAutopilot) CreateContact(contact *model.Contact, authToken string) error {
	requestBody, err := json.Marshal(contact.Data)
	if err != nil {
		log.Printf("Error marshalling json, err: %v", err)
		return err
	}

	request, err := http.NewRequest("POST", autopilotBaseContactURL, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Printf("There is an error while try to create a request for create contact, err: %v", err)
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set(customAutopilotAuthorizationHeader, authToken)

	_, err = a.client.Do(request)
	if err != nil {
		log.Printf("There is an error while try to send a create contact request to Autopilot API, err : %v", err)
		return err
	}
	return nil
}

func (a *apiAutopilot) UpdateContact(contact *model.Contact, authToken string) error {
	requestBody, err := json.Marshal(contact.Data)
	if err != nil {
		log.Printf("Error marshalling json, err: %v", err)
		return err
	}

	request, err := http.NewRequest("POST", autopilotBaseContactURL, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Printf("There is an error while try to create a request for create contact, err: %v", err)
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set(customAutopilotAuthorizationHeader, authToken)

	_, err = a.client.Do(request)
	if err != nil {
		log.Printf("There is an error while try to send a create contact request to Autopilot API, err : %v", err)
		return err
	}
	return nil
}
