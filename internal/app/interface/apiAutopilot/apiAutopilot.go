package apiAutopilot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ramonmacias/autopilot/internal/app/domain/model"
)

type Error struct {
	StatusCode int
	Message    string
}

func (e *Error) Error() string {
	return fmt.Sprintf("Error from autopìlot API with code: %d and body: %s", e.StatusCode, e.Message)
}

type apiAutopilot struct{}

func NewApiAutopilot() *apiAutopilot {
	return &apiAutopilot{}
}

func (a *apiAutopilot) GetContact(id, authToken string) (*model.Contact, error) {
	contactResponse, err := GetContact(id, authToken)
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

func (a *apiAutopilot) CreateContact(contact *model.Contact, authToken string) (*string, error) {
	requestBody, err := json.Marshal(contact.Data)
	if err != nil {
		log.Printf("Error marshalling json, err: %v", err)
		return nil, err
	}

	contactResponse, err := CreateContact(authToken, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Printf("Error sending GET request for get contact information, err: %v", err)
		return nil, err
	}

	return &contactResponse.Id, nil
}

func (a *apiAutopilot) UpdateContact(contact *model.Contact, authToken string) (*string, error) {
	requestBody, err := json.Marshal(contact.Data)
	if err != nil {
		log.Printf("Error marshalling json, err: %v", err)
		return nil, err
	}

	contactResponse, err := UpdateContact(authToken, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Printf("Error sending GET request for get contact information, err: %v", err)
		return nil, err
	}

	return &contactResponse.Id, nil
}
