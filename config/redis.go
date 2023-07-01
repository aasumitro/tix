package config

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
)

func (cfg *Config) InitRedisConn() {
	log.Println("Trying to open redis connection pool . . . .")
	redisSingleton.Do(func() {
		Redis = redis.NewClient(&redis.Options{
			Addr:     cfg.RedisDsnURL,
			Password: cfg.RedisPassword,
			DB:       0,
		})
		if err := Redis.Ping(context.Background()).Err(); err != nil {
			panic(fmt.Sprintf("REDIS_ERROR: %s", err.Error()))
		}
		log.Println("Redis connection pool created . . . .")
	})
}
