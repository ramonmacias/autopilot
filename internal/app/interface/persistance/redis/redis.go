package redis

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/go-redis/redis"
)

var (
	c    *configuration
	once sync.Once
)

type configuration struct {
	client *redis.Client
}

type redisConfigInfo struct {
	RedisDB   string `json:"redis_db"`
	RedisHost string `json:"redis_host"`
	RedisPort string `json:"redis_port"`
	RedisPwd  string `json:"redis_pwd"`
}

func Start() {
	redisConfig()
}

func redisConfig() *configuration {
	once.Do(func() {
		c = &configuration{
			client: redisClient(),
		}
	})
	return c
}

func redisClient() (client *redis.Client) {
	path, err := filepath.Abs("../../config/redis.json")
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
