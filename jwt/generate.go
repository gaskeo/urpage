package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"go-site/errs"
	"go-site/structs"
	"time"
)

const minSecretKeySize = 32

var SecretKey = GenerateKey()

func GenerateJWTToken(id int, expiredAt time.Time, secretKey string) (structs.Payload, string, error) {
	payloadId, err := GenerateId()

	if err != nil {
		return structs.Payload{}, "", err
	}

	payload := structs.Payload{UserId: id,
		PayloadId: payloadId,
		IssuedAt:  time.Now(),
		ExpiredAt: expiredAt}

	if len(secretKey) < minSecretKeySize {
		return structs.Payload{}, "", errs.ErrSmallSecretKey
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &payload)

	JWTString, err := jwtToken.SignedString([]byte(secretKey))

	return payload, JWTString, err
}
