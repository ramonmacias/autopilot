package usecase

import (
	"testing"

	"github.com/ramonmacias/autopilot/internal/app/domain/model"
)

type ContactTestRepository struct {
	contacts map[string]string
}
type ExternalTestApiService struct{}

func (c *ContactTestRepository) FindByID(id string) (*model.Contact, error) {
	_, match := c.contacts[id]
	if !match {
		return nil, nil
	}
	return &model.Contact{
		Id:    id,
		Email: "test@test.com",
		Data:  "Data from cache",
	}, nil
}

func (c *ContactTestRepository) FindByEmail(email string) (*model.Contact, error) {
	_, match := c.contacts[email]
	if !match {
		return nil, nil
	}
	return &model.Contact{
		Id:    "test_id",
		Email: email,
		Data:  "Data from cache",
	}, nil
}

func (c *ContactTestRepository) Save(contact *model.Contact) error {
	c.contacts[contact.Id] = contact.Data
	c.contacts[contact.Email] = contact.Data
	return nil
}

func (c *ContactTestRepository) Delete(contact *model.Contact) error {
	delete(c.contacts, contact.Email)
	delete(c.contacts, contact.Id)
	return nil
}

func (c *ExternalTestApiService) GetContact(id, authToken string) (*model.Contact, error) {
	return &model.Contact{
		Id:    id,
		Email: "test_cached@test.com",
		Data:  "Data external API",
	}, nil
}

func (c *ExternalTestApiService) CreateContact(contact *model.Contact, authToken string) (*string, error) {
	return &contact.Id, nil
}

func (c *ExternalTestApiService) UpdateContact(contact *model.Contact, authToken string) (*string, error) {
	return &contact.Id, nil
}

func TestGetFromExternalAPI(t *testing.T) {
	contacts := make(map[string]string)
	testContactUseCase := NewContactUseCase(&ContactTestRepository{contacts: contacts}, &ExternalTestApiService{})
	contact, _ := testContactUseCase.ShowContact("external_api_id", "auth_test_token")
	if contact.Data != "Data external API" {
		t.Errorf("Expected data (Data external API) but got %s", contact.Data)
	}
}

func TestGetFromCached(t *testing.T) {
	contacts := make(map[string]string)
	testContactUseCase := NewContactUseCase(&ContactTestRepository{contacts: contacts}, &ExternalTestApiService{})
	contact, _ := testContactUseCase.ShowContact("external_api_id", "auth_test_token")
	if contact.Data != "Data external API" {
		t.Errorf("Expected data (Data external API) but got %s", contact.Data)
	}

	contact, _ = testContactUseCase.ShowContact("external_api_id", "auth_test_token")
	if contact.Data != "Data from cache" {
		t.Errorf("Expected data (Data from cache) but got %s", contact.Data)
	}

	contact, _ = testContactUseCase.ShowContact("test_cached@test.com", "auth_test_token")
	if contact.Data != "Data from cache" {
		t.Errorf("Expected data (Data from cache) but got %s", contact.Data)
	}
}

func TestGetFromExternalAPIAfterCreate(t *testing.T) {
	contacts := make(map[string]string)
	testContactUseCase := NewContactUseCase(&ContactTestRepository{contacts: contacts}, &ExternalTestApiService{})
	contact, _ := testContactUseCase.ShowContact("external_api_id", "auth_test_token")
	if contact.Data != "Data external API" {
		t.Errorf("Expected data (Data external API) but got %s", contact.Data)
	}

	contact, _ = testContactUseCase.ShowContact("external_api_id", "auth_test_token")
	if contact.Data != "Data from cache" {
		t.Errorf("Expected data (Data from cache) but got %s", contact.Data)
	}

	testContactUseCase.CreateContact(contact, "auth_test_token")

	contact, _ = testContactUseCase.ShowContact("external_api_id", "auth_test_token")
	if contact.Data != "Data external API" {
		t.Errorf("Expected data (Data external API) but got %s", contact.Data)
	}
}

func TestGetFromExternalAPIAfterUpdate(t *testing.T) {
	contacts := make(map[string]string)
	testContactUseCase := NewContactUseCase(&ContactTestRepository{contacts: contacts}, &ExternalTestApiService{})
	contact, _ := testContactUseCase.ShowContact("external_api_id", "auth_test_token")
	if contact.Data != "Data external API" {
		t.Errorf("Expected data (Data external API) but got %s", contact.Data)
	}

	contact, _ = testContactUseCase.ShowContact("external_api_id", "auth_test_token")
	if contact.Data != "Data from cache" {
		t.Errorf("Expected data (Data from cache) but got %s", contact.Data)
	}

	testContactUseCase.UpdateContact(contact, "auth_test_token")

	contact, _ = testContactUseCase.ShowContact("external_api_id", "auth_test_token")
	if contact.Data != "Data external API" {
		t.Errorf("Expected data (Data external API) but got %s", contact.Data)
	}
}
