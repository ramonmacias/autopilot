package repository

import "github.com/ramonmacias/autopilot/internal/app/domain/model"

type ContactRepository interface {
	FindByID(id string) (*model.Contact, error)
	FindByEmail(email string) (*model.Contact, error)
	Save(*model.Contact) error
	Delete(*model.Contact) error
}
