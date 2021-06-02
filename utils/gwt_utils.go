package utils

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
)

func GenerateId() int64 {

	id, err := rand.Int(rand.Reader, big.NewInt(10000000))
	if err != nil {
		log.Println(err)
	}

	return id.Int64()
}

func GenerateKey() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Println(err)
	}

	return fmt.Sprintf("%x", b)
}
