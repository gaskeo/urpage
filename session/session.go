package session

import (
	"crypto/rand"
	"fmt"
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

func GetCSRFBySessionId(sessionId string) (string, error) {
	CSRFToken, err := redis_api.Get(sessionId)
	return CSRFToken, err
}

func CheckSessionId(writer http.ResponseWriter, request *http.Request) (string, string, error) {
	{ // check cookie
		sessionIdCookie, err := request.Cookie("SessionId")
		if err == nil {
			CSRFToken, err := GetCSRFBySessionId(sessionIdCookie.Value)
			return sessionIdCookie.Value, CSRFToken, err
		}
		sessionId := GenerateSessionId()
		CSRFToken := GenerateCSRFToken()
		expireTime := time.Now().Add(constants.SessionIdExpireTime)
		err = redis_api.SetSession(sessionId, CSRFToken, expireTime)

		AddSessionIdCookie(sessionId, expireTime, writer)

		if err != nil {
			return "", "", err
		}
		return sessionId, CSRFToken, nil
	}
}
