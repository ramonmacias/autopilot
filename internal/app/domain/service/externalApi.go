package service

import "github.com/ramonmacias/autopilot/internal/app/domain/model"

type ExternalApi interface {
	GetContact(id, authToken string) (*model.Contact, error)
	CreateContact(contact *model.Contact, authToken string) error
	UpdateContact(contact *model.Contact, authToken string) error
}
