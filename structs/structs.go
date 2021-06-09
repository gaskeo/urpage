package structs

import (
	"go-site/errs"
	"time"
)

type Data interface{}

type TemplateData map[string]Data

type Answer struct {
	Err string
}

type Payload struct {
	UserId    int
	PayloadId int64
	IssuedAt  time.Time
	ExpiredAt time.Time
}

func (payload *Payload) Valid() error {

	if time.Now().After(payload.ExpiredAt) {
		return errs.ErrExpiredToken
	}

	return nil
}

type User struct {
	UserId     int
	Username   string
	Password   string
	Email      string
	CreateDate time.Time
	ImagePath  string
	Links      [][]string
}
