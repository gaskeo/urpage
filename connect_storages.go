package main

import (
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4"
	"go-site/redis_api"
	"go-site/storage"
	"log"
	"os"
)

func connectStorages() (*pgx.Conn, *redis.Client) {
	var (
		conn *pgx.Conn
		rds  *redis.Client
	)
	var err error

	{ // connect to DB
		conn, err = storage.Connect(os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

		if err != nil {
			log.Fatal("database problem", err)
		}

		log.Println("connected to postgres")
	}

	{ // connect to redis
		rds, err = redis_api.Connect(os.Getenv("REDIS_ADDRESS"), os.Getenv("REDIS_PASSWORD"), 0)

		if err != nil {
			log.Fatal("redis problem", err)
		}

		log.Println("connected to redis")
	}

	return conn, rds
}
