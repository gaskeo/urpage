package session

import (
	"crypto/rand"
	"fmt"
	"go-site/redis_api"
	"log"
	"time"
)

func GenerateSessionId() string {
	sessionId := make([]byte, 32)
	_, err := rand.Read(sessionId)

	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%x", sessionId)
}

func GenerateCSRFToken() string {
	token := make([]byte, 16)
	_, err := rand.Read(token)

	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%x", token)
}

func AddInRedis(sessionId, CSRFToken string, expireTime time.Time) error {
	err := redis_api.Set(sessionId, CSRFToken, expireTime)
	return err
}

func GetCSRFBySessionId(sessionId string) (string, error) {
	CSRFToken, err := redis_api.Get(sessionId)
	return CSRFToken, err
}
