package apiAutopilot

import (
	"fmt"
	"io"
)

const (
	autopilotBaseContactURL = "/contact"
)

type ContactResponse struct {
	Email string `json:"Email"`
	Id    string `json:"contact_id"`
	Body  []byte `json:"-"`
}

func GetContact(id, authToken string) (*ContactResponse, error) {
	return SendRequest("GET", fmt.Sprintf("%s%s/%s", GetBaseUrl(), autopilotBaseContactURL, id), authToken, nil)
}

func UpdateContact(authToken string, body io.Reader) (*ContactResponse, error) {
	return SendRequest("POST", fmt.Sprintf("%s%s", GetBaseUrl(), autopilotBaseContactURL), authToken, body)
}

func CreateContact(authToken string, body io.Reader) (*ContactResponse, error) {
	return SendRequest("POST", fmt.Sprintf("%s%s", GetBaseUrl(), autopilotBaseContactURL), authToken, body)
}
