package tools

import (
	"github.com/go-redis/redis"
)

var RedisClient *redis.Client

func InitRedis(addr, password string, db int) error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	_, err := RedisClient.Ping().Result()
	return err
}
func CloseRedis() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}
