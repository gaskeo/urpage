package main

import (
	"go-site/redis_api"
	"go-site/storage"
	"log"
	"os"
)

func connectStorages() {
	{
		_, err := storage.Connect(os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

		if err != nil {
			log.Fatal("database problem", err)
		}

		log.Println("connected to postgres")
	}

	{
		connectToRedis, err := redis_api.Connect(os.Getenv("REDIS_ADDRESS"), os.Getenv("REDIS_PASSWORD"), 0)

		if err != nil || !connectToRedis {
			log.Fatal("redis problem", err)
		}

		log.Println("connected to redis")
	}
}
