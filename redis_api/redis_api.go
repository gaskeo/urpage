package redis_api

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

var ctx = context.Background()

var rdb *redis.Client

func Connect(address string, password string, db int) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})
}

func Set(key string, value string, expiredDate time.Time) {
	var expiredTime time.Duration

	zeroTime := time.Time{}

	if expiredDate == zeroTime {
		expiredTime = time.Duration(0)
	} else {
		expiredTime = expiredDate.Sub(time.Now())
	}
	fmt.Println(expiredTime)
	fmt.Println(rdb.Set(ctx, key, value, expiredTime))

}
