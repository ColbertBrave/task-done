package redis

import (
	"github.com/go-redis/redis"
	"github.com/task-done/infrastructure/config"
)

var client *redis.Client

func Init() {
	host := config.GetConfig().Redis.Host
	port := config.GetConfig().Redis.Port
	password := config.GetConfig().Redis.Password

	client = redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       0,
	})
}

func Get(key string) (string, error) {
	return client.Get(key).Result()
}
