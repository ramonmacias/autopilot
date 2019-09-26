package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/ramonmacias/autopilot/internal/app/domain/model"
	"github.com/ramonmacias/autopilot/internal/app/domain/repository"
	"github.com/ramonmacias/autopilot/internal/app/interface/apiAutopilot"
	"github.com/ramonmacias/autopilot/internal/app/interface/persistence/redis"
	"github.com/ramonmacias/autopilot/internal/app/usecase"
)

type ContactRequest struct {
	ContactInfo ContactInfoRequest `json:"contact"`
}

type ContactInfoRequest struct {
	Email string `json:"Email"`
}

var (
	client            http.Client
	contactRepository repository.ContactRepository
	contactUseCase    usecase.ContactUseCase
)

const (
	customAutopilotAuthorizationHeader = "autopilotapikey"
)

func init() {
	timeout := time.Duration(5 * time.Second)
	client = http.Client{
		Timeout: timeout,
	}
	contactUseCase = usecase.NewContactUseCase(
		redis.NewContactController(redis.GetClient()),
		apiAutopilot.NewApiAutopilot(client),
	)
}

func ShowContact(w http.ResponseWriter, r *http.Request) {
	contact, err := contactUseCase.ShowContact(mux.Vars(r)["id"], r.Header.Get(customAutopilotAuthorizationHeader))
	w.Header().Set("Content-Type", "application/json")
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
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body request, err: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	requestInfo := &ContactRequest{}
	json.Unmarshal(requestBody, requestInfo)

	err = contactUseCase.UpdateContact(
		model.NewContact("", requestInfo.ContactInfo.Email, string(requestBody)),
		r.Header.Get(customAutopilotAuthorizationHeader),
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
}

func CreateContact(w http.ResponseWriter, r *http.Request) {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body request, err: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	requestInfo := &ContactRequest{}
	json.Unmarshal(requestBody, requestInfo)

	err = contactUseCase.CreateContact(
		model.NewContact("", requestInfo.ContactInfo.Email, string(requestBody)),
		r.Header.Get(customAutopilotAuthorizationHeader),
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
	w.WriteHeader(http.StatusOK)
}
