package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"go-site/constants"
	"go-site/structs"
	"time"
)

var SecretKey = GenerateKey()

func GenerateJWTToken(id int) (structs.Payload, string, time.Time, error) {
	tokenExpireDate := time.Now().Add(constants.JWTExpireTime)

	payloadId, err := GenerateId()

	if err != nil {
		return structs.Payload{}, "", time.Time{}, err
	}

	payload := structs.Payload{UserId: id,
		PayloadId: payloadId,
		IssuedAt:  time.Now(),
		ExpiredAt: tokenExpireDate}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &payload)

	JWTString, err := jwtToken.SignedString([]byte(SecretKey))

	return payload, JWTString, tokenExpireDate, err
}
