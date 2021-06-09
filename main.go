package main

import (
	"log"
	"net/http"
)

func main() {
	conn, rds := connectStorages()

	generateHandlers(conn, rds)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
