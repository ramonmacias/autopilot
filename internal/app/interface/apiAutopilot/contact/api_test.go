package contact

import (
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/ramonmacias/autopilot/internal/app/domain/model"
	"github.com/ramonmacias/autopilot/internal/app/interface/apiAutopilot"
)

type externalAPICLientTest struct{}

func (c *externalAPICLientTest) GetClient() *http.Client {
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	return &client
}

func (c *externalAPICLientTest) GetBaseUrl() string {
	return "https://api2.autopilothq.com/v1"
}

func (c *externalAPICLientTest) SendRequest(method, url, authToken string, body io.Reader) (*apiAutopilot.ContactResponse, error) {
	return &apiAutopilot.ContactResponse{
		Id:    "test_id",
		Email: "test@test.com",
		Body:  []byte("test data"),
	}, nil
}

func NewExternalAPICLientTest() *externalAPICLientTest {
	return &externalAPICLientTest{}
}

func TestGetContact(t *testing.T) {
	api := NewContactAPI(NewExternalAPICLientTest())
	contact, err := api.GetContact("test_id", "test_token")
	if err != nil {
		t.Errorf("Should not be any error on this test but got err: %v", err)
	}
	if contact == nil {
		t.Error("contact shouldn't be nil")
	}
}

func TestCreateContact(t *testing.T) {
	api := NewContactAPI(NewExternalAPICLientTest())
	contact, err := api.CreateContact(model.NewContact("testID", "test@test.com", "test_data"), "test_token")
	if err != nil {
		t.Errorf("Should not be any error on this test but got err: %v", err)
	}
	if contact == nil {
		t.Error("contact shouldn't be nil")
	}
}

func TestUpdateContact(t *testing.T) {
	api := NewContactAPI(NewExternalAPICLientTest())
	contact, err := api.UpdateContact(model.NewContact("testID", "test@test.com", "test_data"), "test_token")
	if err != nil {
		t.Errorf("Should not be any error on this test but got err: %v", err)
	}
	if contact == nil {
		t.Error("contact shouldn't be nil")
	}
}
