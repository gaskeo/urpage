package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"os"
	"time"
)

type User struct {
	userId     int
	username   string
	password   string
	email      string
	createDate time.Time
	imagePath  string
	links      string
}

type Page struct {
	pageId   int
	authorId int
	links    string
}

func connect(username string, password string, dbname string) *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), "postgres://"+username+":"+password+"@localhost:5432/"+dbname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return conn
}

func getUserViaId(conn *pgx.Conn, userId int) *User {
	user := User{}
	err := conn.QueryRow(context.Background(), "SELECT * from user_info where user_id=$1", userId).Scan(
		&user.userId,
		&user.username,
		&user.password,
		&user.email,
		&user.createDate,
		&user.imagePath,
		&user.links)
	if err != nil {
		return &User{}
	}
	return &user
}

func addUser(conn *pgx.Conn, username string, password string, email string, imagePath string, links string) int {
	userId := -1
	err := conn.QueryRow(context.Background(),
		"INSERT INTO user_info (username, password, email, create_date, image_path, links)"+
			"VALUES ($1, $2, $3, $4, $5) RETURNING user_id", username, password, email, time.Now(), imagePath, links).Scan(&userId)
	if err != nil {
		return -1
	}
	return userId
}
