package storage

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"log"
	"strings"
	"time"
	"urpage/constants"
	"urpage/utils"
)

var ErrWrongPassword = errors.New("wrong password")

type User struct {
	UserId     int
	Username   string
	Password   string
	Email      string
	CreateDate time.Time
	ImagePath  string
	Links      [][]string
}

func Connect(host string, username string, password string, dbname string) (*pgx.Conn, error) {
	var (
		conn *pgx.Conn
		err  error
	)

	conn, err = pgx.Connect(context.Background(), "postgres://"+username+":"+password+"@"+host+"/"+dbname)

	return conn, err
}

func GetUserViaId(conn *pgx.Conn, userId int) (User, error) {
	user := User{}

	var (
		imageDB *string
		linksDB *string
		err     error
	)

	{
		err = conn.QueryRow(context.Background(), "SELECT * from user_info WHERE user_id=$1", userId).Scan(
			&user.UserId,
			&user.Username,
			&user.Password,
			&user.Email,
			&user.CreateDate,
			&imageDB,
			&linksDB)

		if err != nil {
			log.Println(err)
			return User{}, err
		}
	}

	{
		if imageDB != nil && len(*imageDB) != 0 {
			user.ImagePath = constants.UserImages + *imageDB
		} else {
			user.ImagePath = constants.UserImages + "default.jpeg"
		}

		if linksDB != nil && len(*linksDB) != 0 {
			linksLst := strings.Split(*linksDB, " ")
			user.Links, err = utils.CreateIconLinkPairs(linksLst)
			if err != nil {
				return User{}, err
			}
		}
	}

	return user, nil
}

func AddUser(conn *pgx.Conn, username string, password string, email string, imagePath string, links string) (int, error) {
	var userId int

	err := conn.QueryRow(context.Background(),
		"INSERT INTO user_info (Username, Password, Email, create_date, image_path, Links)"+
			"VALUES ($1, $2, $3, $4, $5, $6) RETURNING user_id", username, password, email, time.Now(), imagePath, links).Scan(&userId)

	if err != nil {
		return -1, err
	}

	return userId, nil
}

func CheckEmailExistInDB(conn *pgx.Conn, email string) (bool, error) {
	var emailDB string

	err := conn.QueryRow(context.Background(), "SELECT email from user_info WHERE email=$1", email).Scan(&emailDB)

	if err != nil {
		log.Println(err)
		return false, err
	}

	return emailDB == email, nil
}

func GetUserByEmailAndPassword(conn *pgx.Conn, email string, password string) (User, error) {
	user := User{}

	var imageDB *string
	var linksDB *string
	var err error

	{
		err = conn.QueryRow(context.Background(), "SELECT * from user_info WHERE email=$1", email).Scan(
			&user.UserId,
			&user.Username,
			&user.Password,
			&user.Email,
			&user.CreateDate,
			&imageDB,
			&linksDB)

		if err != nil {
			return User{}, err
		}
	}

	{
		PasswordsCompare, err := CheckPassword(password, user.Password)

		if err != nil || !PasswordsCompare {
			return User{}, ErrWrongPassword
		}
	}

	{
		if imageDB != nil && len(*imageDB) != 0 {
			user.ImagePath = constants.UserImages + *imageDB
		} else {
			user.ImagePath = constants.UserImages + "default.jpeg"
		}

		if linksDB != nil && len(*linksDB) != 0 {
			linksLst := strings.Split(*linksDB, " ")
			user.Links, err = utils.CreateIconLinkPairs(linksLst)

			if err != nil {
				return User{}, err
			}

		}
	}

	return user, nil
}

func UpdateUser(conn *pgx.Conn, user User) error {
	var userId int

	links := utils.CreateDBLinksFromPairs(user.Links)

	if strings.Contains(user.ImagePath, constants.UserImages) {
		user.ImagePath = user.ImagePath[len(constants.UserImages):]
	}

	err := conn.QueryRow(context.Background(), "UPDATE user_info SET "+
		"username=$1, email=$2, password=$3, image_path=$4, links=$5 WHERE user_id=$6 RETURNING user_id",
		user.Username, user.Email, user.Password, user.ImagePath, links, user.UserId).Scan(&userId)

	return err
}
