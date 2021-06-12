package main

import (
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4"
	"log"
	"os"
	"urpage/redis_api"
	"urpage/storage"
)

func connectStorages() (*pgx.Conn, *redis.Client) {
	var (
		conn *pgx.Conn
		rdb  *redis.Client
	)
	var err error

	{ // connect to DB
		conn, err = storage.Connect(os.Getenv("DB_ADDRESS"), os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

		if err != nil {
			log.Fatal("database problem", err)
		}

		log.Println("connected to postgres")
	}

	{ // connect to redis
		rdb, err = redis_api.Connect(os.Getenv("REDIS_ADDRESS"), os.Getenv("REDIS_PASSWORD"), 0)

		if err != nil {
			log.Fatal("redis problem", err)
		}

		log.Println("connected to redis")
	}

	return conn, rdb
}
