package infra

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

func InitRedisClient(host, port, password string, db int) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:        fmt.Sprintf("%s:%s", host, port),
		Password:    password,
		DB:          db,
		DialTimeout: 3 * time.Second,
	})
	result := rdb.Ping(context.TODO())
	return rdb, result.Err()
}
