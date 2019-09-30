package contact

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ramonmacias/autopilot/internal/app/domain/model"
	"github.com/ramonmacias/autopilot/internal/app/interface/apiAutopilot"
)

const (
	autopilotBaseContactURL = "/contact"
)

type contactApi struct {
	client apiAutopilot.APIClientable
}

func NewContactAPI(client apiAutopilot.APIClientable) *contactApi {
	return &contactApi{
		client: client,
	}
}

func (a *contactApi) GetContact(id, authToken string) (*model.Contact, error) {
	contactResponse, err := a.client.SendRequest("GET", fmt.Sprintf("%s%s/%s", a.client.GetBaseUrl(), autopilotBaseContactURL, id), authToken, nil)
	if err != nil {
		log.Printf("Error sending GET request for get contact information, err: %v", err)
		return nil, err
	}

	return &model.Contact{
		Id:    contactResponse.Id,
		Email: contactResponse.Email,
		Data:  string(contactResponse.Body),
	}, nil
}

func (a *contactApi) CreateContact(contact *model.Contact, authToken string) (*string, error) {
	requestBody, err := json.Marshal(contact.Data)
	if err != nil {
		log.Printf("Error marshalling json, err: %v", err)
		return nil, err
	}

	contactResponse, err := a.client.SendRequest("POST", fmt.Sprintf("%s%s", a.client.GetBaseUrl(), autopilotBaseContactURL), authToken, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Printf("Error sending GET request for get contact information, err: %v", err)
		return nil, err
	}

	return &contactResponse.Id, nil
}

func (a *contactApi) UpdateContact(contact *model.Contact, authToken string) (*string, error) {
	requestBody, err := json.Marshal(contact.Data)
	if err != nil {
		log.Printf("Error marshalling json, err: %v", err)
		return nil, err
	}

	contactResponse, err := a.client.SendRequest("POST", fmt.Sprintf("%s%s", a.client.GetBaseUrl(), autopilotBaseContactURL), authToken, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Printf("Error sending GET request for get contact information, err: %v", err)
		return nil, err
	}

	return &contactResponse.Id, nil
}
