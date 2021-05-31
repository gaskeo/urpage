package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"go-site/constants"
	"go-site/utils"
	"log"
	"os"
	"strings"
	"time"
)

var conn *pgx.Conn

type User struct {
	UserId     int
	Username   string
	Password   string
	Email      string
	CreateDate time.Time
	ImagePath  string
	Links      [][]string
}

func Connect(username string, password string, dbname string) *pgx.Conn {
	connect, err := pgx.Connect(context.Background(), "postgres://"+username+":"+password+"@localhost:5432/"+dbname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	conn = connect
	return conn
}

func GetUserViaId(userId int) *User {
	user := User{}

	var image *string
	var links *string

	err := conn.QueryRow(context.Background(), "SELECT * from user_info WHERE user_id=$1", userId).Scan(
		&user.UserId,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.CreateDate,
		&image,
		&links)
	if err != nil {
		log.Println(err)
		return &User{}
	}

	if image != nil && len(*image) != 0 {
		user.ImagePath = constants.UserImages + *image
	} else {
		user.ImagePath = constants.UserImages + "test.jpeg"
	}

	if links != nil && len(*links) != 0 {
		linksLst := strings.Split(*links, " ")
		user.Links = utils.CreateIconLinkPairs(linksLst)
	}

	return &user
}

func AddUser(username string, password string, email string, imagePath string, links string) int {
	userId := -1
	err := conn.QueryRow(context.Background(),
		"INSERT INTO user_info (Username, Password, Email, create_date, image_path, Links)"+
			"VALUES ($1, $2, $3, $4, $5, $6) RETURNING user_id", username, password, email, time.Now(), imagePath, links).Scan(&userId)
	if err != nil {
		log.Println(err)
		return -1
	}
	return userId
}

func CheckEmailExistInDB(email string) bool {
	var emailInDB string

	err := conn.QueryRow(context.Background(), "SELECT email from user_info WHERE email=$1", email).Scan(&emailInDB)
	if err != nil {
		log.Println(err)
	}

	return emailInDB == email
}

func GetUserByEmailAndPassword(email string, password string) *User {
	user := User{}

	var image *string
	var links *string

	err := conn.QueryRow(context.Background(), "SELECT * from user_info WHERE email=$1", email).Scan(
		&user.UserId,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.CreateDate,
		&image,
		&links)
	if err != nil {
		log.Println(err)
		return &User{}
	}

	if password == user.Password {
		return &User{}
	}

	if image != nil && len(*image) != 0 {
		user.ImagePath = constants.UserImages + *image
	} else {
		user.ImagePath = constants.UserImages + "test.jpeg"
	}

	if links != nil && len(*links) != 0 {
		linksLst := strings.Split(*links, " ")
		user.Links = utils.CreateIconLinkPairs(linksLst)
	}

	return &user
}
