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
		&user.imagePath)
	if err != nil {
		return &User{}
	}
	return &user
}

func addUser(conn *pgx.Conn, username string, password string, email string, imagePath string) int {
	userId := -1
	err := conn.QueryRow(context.Background(),
		"INSERT INTO user_info (username, password, email, create_date, image_path)"+
			"VALUES ($1, $2, $3, $4) RETURNING user_id", username, password, email, time.Now(), imagePath).Scan(&userId)
	if err != nil {
		return -1
	}
	return userId
}

func getPageViaId(conn *pgx.Conn, pageId int) *Page {
	page := Page{}
	err := conn.QueryRow(context.Background(),
		"SELECT * FROM page_info WHERE page_id == $1", pageId).Scan(&page.pageId, &page.authorId, &page.links)
	if err != nil {
		return &Page{}
	}
	return &page
}

func addPage(conn *pgx.Conn, authorId int, links string) int {
	pageId := -1
	user := getUserViaId(conn, authorId)
	if user.userId == 0 {
		return -1
	}
	err := conn.QueryRow(context.Background(),
		"INSERT INTO page_info (author_id, links) VALUES ($1, $2)", authorId, links).Scan(&pageId)
	if err != nil {
		return -1
	}
	return pageId
}
