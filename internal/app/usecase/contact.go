package usecase

import (
	"github.com/ramonmacias/autopilot/internal/app/domain/model"
	"github.com/ramonmacias/autopilot/internal/app/domain/repository"
	"github.com/ramonmacias/autopilot/internal/app/domain/service"
)

type ContactUseCase interface {
	ShowContact(id, authToken string) (*model.Contact, error)
	CreateContact(contact *model.Contact, authToken string) error
	UpdateContact(contact *model.Contact, authToken string) error
}

type contactUseCase struct {
	repo       repository.ContactRepository
	httpClient service.ExternalApi
}

func NewContactUseCase(repo repository.ContactRepository, httpClient service.ExternalApi) *contactUseCase {
	return &contactUseCase{
		repo:       repo,
		httpClient: httpClient,
	}
}

func (c *contactUseCase) ShowContact(id string, authToken string) (*model.Contact, error) {
	contact, err := c.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if contact == nil {
		contact, err = c.repo.FindByEmail(id)
		if err != nil {
			return nil, err
		}
	}
	if contact == nil {
		contact, err = c.httpClient.GetContact(id, authToken)
		if err != nil {
			return nil, err
		}
		if err = c.repo.Save(contact); err != nil {
			return nil, err
		}
	}
	return contact, nil
}

func (c *contactUseCase) CreateContact(contact *model.Contact, authToken string) error {
	id, err := c.httpClient.CreateContact(contact, authToken)
	if err != nil {
		return err
	}
	contact.Id = *id
	return c.repo.Delete(contact)
}

func (c *contactUseCase) UpdateContact(contact *model.Contact, authToken string) error {
	id, err := c.httpClient.UpdateContact(contact, authToken)
	if err != nil {
		return err
	}
	contact.Id = *id
	return c.repo.Delete(contact)
}
