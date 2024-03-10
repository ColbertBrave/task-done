package redis

import (
	"github.com/cloud-disk/infrastructure/config"
	"github.com/go-redis/redis"
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
