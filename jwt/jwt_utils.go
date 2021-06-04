package jwt

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
)

func GenerateId() (int64, error) {
	id, err := rand.Int(rand.Reader, big.NewInt(10000000))

	if err != nil {
		return 0, err
	}

	return id.Int64(), nil
}

func GenerateKey() string {
	key := make([]byte, 32)
	_, err := rand.Read(key)

	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%x", key)
}

func GenerateRefreshToken() (string, error) {
	refreshToken := make([]byte, 16)
	_, err := rand.Read(refreshToken)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", refreshToken), nil
}
