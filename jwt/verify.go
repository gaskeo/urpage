package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"go-site/errs"
	"go-site/structs"
)

func VerifyToken(token string, secretKey string) (*structs.Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {

		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errs.ErrInvalidToken
		}

		return []byte(secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &structs.Payload{}, keyFunc)

	if err != nil {
		ver, ok := err.(*jwt.ValidationError)

		if ok && errors.Is(ver.Inner, errs.ErrExpiredToken) {
			return nil, errs.ErrExpiredToken
		}

		return nil, errs.ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*structs.Payload)

	if !ok {
		return nil, errs.ErrInvalidToken
	}

	return payload, nil
}
