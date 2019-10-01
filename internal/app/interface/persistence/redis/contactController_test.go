package redis

import (
	"fmt"
	"testing"

	"github.com/go-redis/redis"
	"github.com/ramonmacias/autopilot/internal/app/domain/model"
)

var (
	client *redis.Client
)

func init() {
	client = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf(
			"%s:%s",
			"127.0.0.1",
			"6379",
		),
		Password: "",
		DB:       0,
	})
}

func TestFindByID(t *testing.T) {
	if err := client.FlushAll().Err(); err != nil {
		t.Errorf("There is an error with your test redis server: %v", err)
	}
	contactController := NewContactController(client)
	contact, err := contactController.FindByID("test_id")
	if err != nil || contact != nil {
		t.Error("When a key doesn't exists we need to return nil in both results but got some of them not nil")
	}

	contactController.Save(model.NewContact("test_id", "email", "test_data"))

	contact, err = contactController.FindByID("test_id")
	if err != nil {
		t.Errorf("Expect not error but got %v", err)
	}
	if contact == nil {
		t.Error("Expected contact object not nil")
	}
	if contact.Id != "test_id" {
		t.Errorf("Expected id test_id but got %s", contact.Id)
	}
	if contact.Data != "test_data" {
		t.Errorf("Expected test_data but got %s", contact.Data)
	}
}

func TestFindByEmail(t *testing.T) {
	if err := client.FlushAll().Err(); err != nil {
		t.Errorf("There is an error with your test redis server: %v", err)
	}
	contactController := NewContactController(client)
	contact, err := contactController.FindByEmail("test@test.com")
	if err != nil || contact != nil {
		t.Error("When a key doesn't exists we need to return nil in both results but got some of them not nil")
	}

	contactController.Save(model.NewContact("test_id", "test@test.com", "test_data"))

	contact, err = contactController.FindByEmail("test@test.com")
	if err != nil {
		t.Errorf("Expect not error but got %v", err)
	}
	if contact == nil {
		t.Error("Expected contact object not nil")
	}
	if contact.Email != "test@test.com" {
		t.Errorf("Expected email test@test.com but got %s", contact.Email)
	}
	if contact.Data != "test_data" {
		t.Errorf("Expected test_data but got %s", contact.Data)
	}
}

func TestSaveAndDeleteMethods(t *testing.T) {
	if err := client.FlushAll().Err(); err != nil {
		t.Errorf("There is an error with your test redis server: %v", err)
	}
	contactController := NewContactController(client)

	contact, err := contactController.FindByEmail("test@test.com")
	if err != nil || contact != nil {
		t.Error("When a key doesn't exists we need to return nil in both results but got some of them not nil")
	}

	if err := contactController.Save(model.NewContact("test_id", "test@test.com", "test_data")); err != nil {
		t.Errorf("Expected no error while save but got %v", err)
	}
}
