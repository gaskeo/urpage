package main

import (
	"log"
	"net/http"
)

func main() {
	connectStorages()

	generateHandlers()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
