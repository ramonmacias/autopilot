package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ramonmacias/autopilot/internal/app/domain/model"
	"github.com/ramonmacias/autopilot/internal/app/domain/repository"
	"github.com/ramonmacias/autopilot/internal/app/interface/apiAutopilot"
	"github.com/ramonmacias/autopilot/internal/app/interface/apiAutopilot/contact"
	"github.com/ramonmacias/autopilot/internal/app/interface/persistence/redis"
	"github.com/ramonmacias/autopilot/internal/app/usecase"
)

type ContactRequest struct {
	ContactInfo ContactInfoRequest `json:"contact"`
}

type ContactInfoRequest struct {
	Email string `json:"Email"`
}

type Response struct {
	ContactId string `json:"contact_id"`
}

var (
	client                   http.Client
	contactRepository        repository.ContactRepository
	contactUseCase           usecase.ContactUseCase
	unauthorizedErrorMessage = `{"error": "Bad Request", "message": "No autopilotapikey header, or apikey query parameter provided."}`
)

const (
	customAutopilotAuthorizationHeader = "autopilotapikey"
)

func init() {
	contactUseCase = usecase.NewContactUseCase(
		redis.NewContactController(redis.GetClient()),
		contact.NewContactAPI(apiAutopilot.NewApiClient()),
	)
}

func ShowContact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	authToken := r.Header.Get(customAutopilotAuthorizationHeader)
	if authToken == "" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(unauthorizedErrorMessage))
		return
	}
	contact, err := contactUseCase.ShowContact(mux.Vars(r)["id"], authToken)
	if err != nil {
		switch err.(type) {
		case *apiAutopilot.Error:
			apiError := err.(*apiAutopilot.Error)
			w.WriteHeader(apiError.StatusCode)
			w.Write([]byte(apiError.Message))
		default:
			log.Printf("Unexpected error getting contact: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(contact.Data))
}

func UpdateContact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	authToken := r.Header.Get(customAutopilotAuthorizationHeader)
	if authToken == "" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(unauthorizedErrorMessage))
		return
	}
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body request, err: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	requestInfo := &ContactRequest{}
	json.Unmarshal(requestBody, requestInfo)

	contactId, err := contactUseCase.UpdateContact(
		model.NewContact("", requestInfo.ContactInfo.Email, string(requestBody)),
		authToken,
	)
	if err != nil {
		switch err.(type) {
		case *apiAutopilot.Error:
			apiError := err.(*apiAutopilot.Error)
			w.WriteHeader(apiError.StatusCode)
			w.Write([]byte(apiError.Message))
		default:
			log.Printf("Unexpected error updating contact: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	response := Response{
		ContactId: *contactId,
	}
	json.NewEncoder(w).Encode(&response)
}

func CreateContact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	authToken := r.Header.Get(customAutopilotAuthorizationHeader)
	if authToken == "" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(unauthorizedErrorMessage))
		return
	}

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body request, err: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	requestInfo := &ContactRequest{}
	json.Unmarshal(requestBody, requestInfo)

	contactId, err := contactUseCase.CreateContact(
		model.NewContact("", requestInfo.ContactInfo.Email, string(requestBody)),
		authToken,
	)
	if err != nil {
		switch err.(type) {
		case *apiAutopilot.Error:
			apiError := err.(*apiAutopilot.Error)
			w.WriteHeader(apiError.StatusCode)
			w.Write([]byte(apiError.Message))
		default:
			log.Printf("Unexpected error creating contact: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	response := Response{
		ContactId: *contactId,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&response)
}
