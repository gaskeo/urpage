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
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		log.Println(err)
	}

	return fmt.Sprintf("%x", key)
}

func GenerateRefreshToken() string {
	refreshToken := make([]byte, 16)
	_, err := rand.Read(refreshToken)
	if err != nil {
		log.Println(err)
	}

	return fmt.Sprintf("%x", refreshToken)
}
