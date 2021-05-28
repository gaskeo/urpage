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
	err := conn.QueryRow(context.Background(), "SELECT user_id, username, password, email, create_date from user_info where user_id=$1", userId).Scan(
		&user.userId,
		&user.username,
		&user.password,
		&user.email,
		&user.createDate)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	return &user
}
