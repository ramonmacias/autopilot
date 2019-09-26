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

type ContactResponse struct {
	Email string `json:"Email"`
	Id    string `json:"contact_id"`
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
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contact.Data)

	// contactRepository = redis.NewContactController(redis.GetClient())
	// contact, err := contactRepository.FindByID(mux.Vars(r)["id"])
	// if err != nil {
	// 	log.Printf("Error retrieving a contact, err: %v", err)
	// 	w.WriteHeader(http.StatusBadGateway)
	// 	return
	// }
	// if contact != nil {
	// 	log.Println("FOUND IN CACHE!")
	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.WriteHeader(http.StatusOK)
	// 	json.NewEncoder(w).Encode(contact.Data)
	// } else {
	// 	request, err := http.NewRequest("GET", autopilotBaseContactURL+"/"+mux.Vars(r)["id"], r.Body)
	// 	if err != nil {
	// 		log.Printf("There is an error while try to create a request for create contact, err: %v", err)
	// 		w.WriteHeader(http.StatusBadGateway)
	// 		return
	// 	}
	// 	request.Header.Set(customAutopilotAuthorizationHeader, r.Header.Get(customAutopilotAuthorizationHeader))
	//
	// 	resp, err := client.Do(request)
	// 	if err != nil {
	// 		log.Printf("There is an error while try to send a create contact request to Autopilot API, err : %v", err)
	// 		w.WriteHeader(http.StatusBadGateway)
	// 		return
	// 	}
	// 	body, err := ioutil.ReadAll(resp.Body)
	// 	if err != nil {
	// 		log.Printf("Unexpected error: %v", err)
	// 		w.WriteHeader(http.StatusBadGateway)
	// 		return
	// 	}
	// 	resp.Body.Close()
	// 	resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	//
	// 	contactResponse := &ContactResponse{}
	// 	json.NewDecoder(resp.Body).Decode(contactResponse)
	// 	log.Printf("Contact response: %v", contactResponse)
	//
	// 	if err = contactRepository.Save(model.NewContact(contactResponse.Id, contactResponse.Email, string(body))); err != nil {
	// 		log.Printf("Error saving the contact into the cache, err: %v", err)
	// 		w.WriteHeader(http.StatusBadGateway)
	// 		return
	// 	}
	//
	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.WriteHeader(http.StatusOK)
	// 	w.Write(body)
	// }
}

func UpdateContact(w http.ResponseWriter, r *http.Request) {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body request, err: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	err = contactUseCase.UpdateContact(
		model.NewContact("", "", string(requestBody)),
		r.Header.Get(customAutopilotAuthorizationHeader),
	)
	if err != nil {
		log.Printf("Error updating contact: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
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
	err = contactUseCase.CreateContact(
		model.NewContact("", "", string(requestBody)),
		r.Header.Get(customAutopilotAuthorizationHeader),
	)
	if err != nil {
		log.Printf("Error updating contact: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// requestBody, _ := ioutil.ReadAll(r.Body)
// r.Body.Close()
//
// contactRepository = redis.NewContactController(redis.GetClient())
// contactResponse := &ContactResponse{}
// json.NewDecoder(r.Body).Decode(contactResponse)
// log.Printf("Contact response: %v", contactResponse)
//
// r.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))
//
// request, err := http.NewRequest("POST", autopilotBaseContactURL, r.Body)
// if err != nil {
// 	log.Printf("There is an error while try to create a request for create contact, err: %v", err)
// 	w.WriteHeader(http.StatusBadGateway)
// 	return
// }
// request.Header.Set("Content-Type", "application/json")
// request.Header.Set(customAutopilotAuthorizationHeader, r.Header.Get(customAutopilotAuthorizationHeader))
//
// resp, err := client.Do(request)
// if err != nil {
// 	log.Printf("There is an error while try to send a create contact request to Autopilot API, err : %v", err)
// 	w.WriteHeader(http.StatusBadGateway)
// 	return
// }
// defer resp.Body.Close()
//
// body, err := ioutil.ReadAll(resp.Body)
// if err != nil {
// 	log.Printf("Unexpected error: %v", err)
// 	w.WriteHeader(http.StatusBadGateway)
// 	return
// }
//
// if err := contactRepository.Delete(model.NewContact(contactResponse.Id, contactResponse.Email, "")); err != nil {
// 	log.Printf("Error while try to remove data from cache: %v", err)
// 	w.WriteHeader(http.StatusBadGateway)
// 	return
// }
//
// w.Header().Set("Content-Type", "application/json")
// w.WriteHeader(http.StatusOK)
// w.Write(body)
