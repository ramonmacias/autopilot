package redis

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/ramonmacias/autopilot/internal/app/domain/model"
)

type contactController struct {
	redisClient *redis.Client
}

func NewContactController(client *redis.Client) *contactController {
	return &contactController{
		redisClient: client,
	}
}

func (r contactController) FindByID(id string) (*model.Contact, error) {

	data, err := r.redisClient.Get(id).Result()
	if err != nil {
		return nil, err
	}
	return &model.Contact{
		Id:   id,
		Data: data,
	}, nil
}

func (r contactController) FindByEmail(email string) (*model.Contact, error) {
	data, err := r.redisClient.Get(email).Result()
	if err != nil {
		return nil, err
	}
	return &model.Contact{
		Email: email,
		Data:  data,
	}, nil
}

func (r contactController) Save(contact *model.Contact) error {
	if err := r.redisClient.Set(contact.Id, contact.Data, time.Duration(0)).Err(); err != nil {
		return err
	}
	if err := r.redisClient.Set(contact.Email, contact.Data, time.Duration(0)).Err(); err != nil {
		return err
	}
	return nil
}

func (r contactController) Delete(contact *model.Contact) error {
	if err := r.redisClient.Del(contact.Id, contact.Email).Err(); err != nil {
		return err
	}
	return nil
}
