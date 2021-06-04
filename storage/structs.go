package storage

import "time"

type User struct {
	UserId     int
	Username   string
	Password   string
	Email      string
	CreateDate time.Time
	ImagePath  string
	Links      [][]string
}

type SomeUsers map[string]User
