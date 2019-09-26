package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/ramonmacias/autopilot/internal/app/domain/model"
	"github.com/ramonmacias/autopilot/internal/app/domain/repository"
	"github.com/ramonmacias/autopilot/internal/app/interface/persistance/redis"
)

type ContactResponse struct {
	Email string `json:"Email"`
	Id    string `json:"contact_id"`
}

var (
	client            http.Client
	contactRepository repository.ContactRepository
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
	contactRepository = redis.NewContactController(redis.GetClient())
	contact, err := contactRepository.FindByID(mux.Vars(r)["id"])
	if err != nil {
		log.Printf("Error retrieving a contact, err: %v", err)
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	if contact != nil {
		log.Println("FOUND IN CACHE!")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(contact.Data)
	} else {
		request, err := http.NewRequest("GET", autopilotBaseContactURL+"/"+mux.Vars(r)["id"], r.Body)
		if err != nil {
			log.Printf("There is an error while try to create a request for create contact, err: %v", err)
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		request.Header.Set(customAutopilotAuthorizationHeader, r.Header.Get(customAutopilotAuthorizationHeader))

		resp, err := client.Do(request)
		if err != nil {
			log.Printf("There is an error while try to send a create contact request to Autopilot API, err : %v", err)
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Unexpected error: %v", err)
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		resp.Body.Close()
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		contactResponse := &ContactResponse{}
		json.NewDecoder(resp.Body).Decode(contactResponse)
		log.Printf("Contact response: %v", contactResponse)

		if err = contactRepository.Save(model.NewContact(contactResponse.Id, contactResponse.Email, string(body))); err != nil {
			log.Printf("Error saving the contact into the cache, err: %v", err)
			w.WriteHeader(http.StatusBadGateway)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}
}

func UpdateContact(w http.ResponseWriter, r *http.Request) {
	requestBody, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()

	contactRepository = redis.NewContactController(redis.GetClient())
	contactResponse := &ContactResponse{}
	json.NewDecoder(r.Body).Decode(contactResponse)
	log.Printf("Contact response: %v", contactResponse)

	r.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))

	request, err := http.NewRequest("POST", autopilotBaseContactURL, r.Body)
	if err != nil {
		log.Printf("There is an error while try to create a request for create contact, err: %v", err)
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set(customAutopilotAuthorizationHeader, r.Header.Get(customAutopilotAuthorizationHeader))

	resp, err := client.Do(request)
	if err != nil {
		log.Printf("There is an error while try to send a create contact request to Autopilot API, err : %v", err)
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Unexpected error: %v", err)
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	if err := contactRepository.Delete(model.NewContact(contactResponse.Id, contactResponse.Email, "")); err != nil {
		log.Printf("Error while try to remove data from cache: %v", err)
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func CreateContact(w http.ResponseWriter, r *http.Request) {
	requestBody, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()

	contactRepository = redis.NewContactController(redis.GetClient())
	contactResponse := &ContactResponse{}
	json.NewDecoder(r.Body).Decode(contactResponse)
	log.Printf("Contact response: %v", contactResponse)

	r.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))

	request, err := http.NewRequest("POST", autopilotBaseContactURL, r.Body)
	if err != nil {
		log.Printf("There is an error while try to create a request for create contact, err: %v", err)
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set(customAutopilotAuthorizationHeader, r.Header.Get(customAutopilotAuthorizationHeader))

	resp, err := client.Do(request)
	if err != nil {
		log.Printf("There is an error while try to send a create contact request to Autopilot API, err : %v", err)
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Unexpected error: %v", err)
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	if err := contactRepository.Delete(model.NewContact(contactResponse.Id, contactResponse.Email, "")); err != nil {
		log.Printf("Error while try to remove data from cache: %v", err)
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
