package session

import (
	"crypto/rand"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go-site/constants"
	"go-site/redis_api"
	"log"
	"net/http"
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

func GetCSRFBySessionId(rds *redis.Client, sessionId string) (string, error) {
	CSRFToken, err := redis_api.Get(rds, sessionId)
	return CSRFToken, err
}

func CheckSessionId(writer http.ResponseWriter, request *http.Request, rds *redis.Client) (string, string, error) {
	{ // check cookie
		sessionIdCookie, err := request.Cookie("SessionId")
		if err == nil {
			CSRFToken, err := GetCSRFBySessionId(rds, sessionIdCookie.Value)
			return sessionIdCookie.Value, CSRFToken, err
		}
		sessionId := GenerateSessionId()
		CSRFToken := GenerateCSRFToken()
		expireTime := time.Now().Add(constants.SessionIdExpireTime)
		err = redis_api.SetSession(rds, sessionId, CSRFToken, expireTime)

		AddSessionIdCookie(sessionId, expireTime, writer)

		if err != nil {
			return "", "", err
		}
		return sessionId, CSRFToken, nil
	}
}

func DeleteSession(writer http.ResponseWriter, rds *redis.Client, sessionId string) error {
	DeleteSessionIdCookie(writer)
	return redis_api.DeleteSession(rds, sessionId)
}
