package main

import (
	"log"
	"net/http"
)

func main() {
	conn, rdb := connectStorages()

	generateHandlers(conn, rdb)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
