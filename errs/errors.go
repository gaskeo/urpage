package errs

import "errors"

// jwt_api errors
var ErrExpiredToken = errors.New("token has expired")
var ErrSmallSecretKey = errors.New("small secret key")
var ErrInvalidToken = errors.New("invalid token")
var ErrInvalidRefreshToken = errors.New("invalid refresh token")

// DB errors
var ErrWrongPassword = errors.New("wrong password")
