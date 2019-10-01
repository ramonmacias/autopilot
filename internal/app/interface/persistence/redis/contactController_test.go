package redis

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/go-redis/redis"
	"github.com/ramonmacias/autopilot/internal/app/domain/model"
)

var (
	client     *redis.Client
	controller *contactController
)

func connectingToRedis() (client *redis.Client) {
	path, err := filepath.Abs("../../../../../config/test/redis.json")
	if err != nil {
		log.Printf("Error while try to get abs path: %v", err)
	}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panicf("There is an error while try to read redis config file: %v", err)
	}
	redisInfo := &redisConfigInfo{}
	if err = json.Unmarshal([]byte(file), redisInfo); err != nil {
		log.Panicf("There is an error while try to unmarshal the json redis config info, err: %v", err)
	}
	redisDBNum, err := strconv.Atoi(redisInfo.RedisDB)
	if err != nil {
		log.Panicf("Incorrect format for json field redis_db, err: %v", err)
	}
	client = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf(
			"%s:%s",
			redisInfo.RedisHost,
			redisInfo.RedisPort,
		),
		Password: redisInfo.RedisPwd,
		DB:       redisDBNum,
	})
	_, err = client.Ping().Result()
	if err == nil {
		log.Println("Redis initialized successfully")
		return client
	} else {
		log.Printf("Error connecting to Redis, err: %v", err)
	}
	return client
}

func init() {
	client = connectingToRedis()
	controller = NewContactController(client)
}

func TestFindByID(t *testing.T) {
	if err := client.FlushAll().Err(); err != nil {
		t.Errorf("There is an error with your test redis server: %v", err)
	}
	contact, err := controller.FindByID("test_id")
	if err != nil || contact != nil {
		t.Error("When a key doesn't exists we need to return nil in both results but got some of them not nil")
	}

	controller.Save(model.NewContact("test_id", "email", "test_data"))

	contact, err = controller.FindByID("test_id")
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
	contact, err := controller.FindByEmail("test@test.com")
	if err != nil || contact != nil {
		t.Error("When a key doesn't exists we need to return nil in both results but got some of them not nil")
	}

	controller.Save(model.NewContact("test_id", "test@test.com", "test_data"))

	contact, err = controller.FindByEmail("test@test.com")
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

	contact, err := controller.FindByEmail("test@test.com")
	if err != nil || contact != nil {
		t.Error("When a key doesn't exists we need to return nil in both results but got some of them not nil")
	}

	if err := controller.Save(model.NewContact("test_id", "test@test.com", "test_data")); err != nil {
		t.Errorf("Expected no error while save but got %v", err)
	}

	contact, err = controller.FindByEmail("test@test.com")
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if contact == nil {
		t.Error("Expected not nil contact but got nil")
	}

	if err := controller.Delete(model.NewContact("test_id", "test@test.com", "test_data")); err != nil {
		t.Errorf("Expected no error while delete but got %v", err)
	}

	contact, err = controller.FindByEmail("test@test.com")
	if err != nil || contact != nil {
		t.Error("When a key doesn't exists we need to return nil in both results but got some of them not nil")
	}

}
