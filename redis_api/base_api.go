package redis_api

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

var ctx = context.Background()

func Connect(address string, password string, db int) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})

	_, err := rdb.Ping(ctx).Result()

	return rdb, err
}

func Set(rdb *redis.Client, key string, value string, expiredDate time.Time) error {
	var expiredTime time.Duration

	zeroTime := time.Time{}

	if expiredDate == zeroTime {
		expiredTime = time.Duration(0)
	} else {
		expiredTime = expiredDate.Sub(time.Now())
	}

	_, err := rdb.Set(ctx, key, value, expiredTime).Result()

	return err
}

func Get(rdb *redis.Client, key string) (string, error) {
	result, err := rdb.Get(ctx, key).Result()

	if err != nil {
		return "", err
	}

	return result, nil
}
