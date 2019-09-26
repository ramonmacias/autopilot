package apiAutopilot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ramonmacias/autopilot/internal/app/domain/model"
)

const (
	customAutopilotAuthorizationHeader = "autopilotapikey"
	autopilotBaseContactURL            = "https://api2.autopilothq.com/v1/contact"
)

type Error struct {
	StatusCode int
	Message    string
}

func (e *Error) Error() string {
	return fmt.Sprintf("Error from autopìlot API with code: %d and body: %s", e.StatusCode, e.Message)
}

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
	body, _ := ioutil.ReadAll(resp.Body)
	if err != nil || resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &Error{
			StatusCode: resp.StatusCode,
			Message:    string(body),
		}
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

func (a *apiAutopilot) CreateContact(contact *model.Contact, authToken string) (*string, error) {
	requestBody, err := json.Marshal(contact.Data)
	if err != nil {
		log.Printf("Error marshalling json, err: %v", err)
		return nil, err
	}

	request, err := http.NewRequest("POST", autopilotBaseContactURL, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Printf("There is an error while try to create a request for create contact, err: %v", err)
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set(customAutopilotAuthorizationHeader, authToken)

	resp, err := a.client.Do(request)
	body, _ := ioutil.ReadAll(resp.Body)
	if err != nil || resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &Error{
			StatusCode: resp.StatusCode,
			Message:    string(body),
		}
	}

	response := &ContactResponse{}
	json.Unmarshal(body, response)
	return &response.Id, nil
}

func (a *apiAutopilot) UpdateContact(contact *model.Contact, authToken string) (*string, error) {
	requestBody, err := json.Marshal(contact.Data)
	if err != nil {
		log.Printf("Error marshalling json, err: %v", err)
		return nil, err
	}

	request, err := http.NewRequest("POST", autopilotBaseContactURL, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Printf("There is an error while try to create a request for create contact, err: %v", err)
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set(customAutopilotAuthorizationHeader, authToken)

	resp, err := a.client.Do(request)
	body, _ := ioutil.ReadAll(resp.Body)
	if err != nil || resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &Error{
			StatusCode: resp.StatusCode,
			Message:    string(body),
		}
	}

	response := &ContactResponse{}
	json.Unmarshal(body, response)
	return &response.Id, nil
}
