package gateway

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(host, port, password string) (*redis.Client, error) {

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       0, // use default DB
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return rdb, nil
}

func GenKey(tagId, featureId int32) string {
	return fmt.Sprintf("t%d f%d", tagId, featureId)
}
