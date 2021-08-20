package jwt_api

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"os"
)

func GenerateId() (int64, error) {
	id, err := rand.Int(rand.Reader, big.NewInt(10000000))

	if err != nil {
		return 0, err
	}

	return id.Int64(), nil
}

func GenerateKey() string {
	if os.Getenv("DEBUG") == "1" {
		return "1"
	}

	key := make([]byte, 32)
	_, err := rand.Read(key)

	if err != nil {
		log.Fatal(err)
		return ""
	}

	return fmt.Sprintf("%x", key)
}
