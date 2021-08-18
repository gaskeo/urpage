package storage

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"log"
	"time"
	"urpage/constants"
)

var ErrWrongPassword = errors.New("wrong password")

type User struct {
	UserId      int
	Username    string
	Description string
	Password    string
	Email       string
	CreateDate  time.Time
	ImagePath   string
	Links       [][]string
	LikesCount  int
	Verified    bool
}

type LinkWithPath struct {
	Link string
	Path string
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
		rows pgx.Rows
		err  error
	)

	{
		err = conn.QueryRow(context.Background(), "SELECT * from users WHERE user_id=$1", userId).Scan(
			&user.UserId,
			&user.Username,
			&user.Description,
			&user.Password,
			&user.Email,
			&user.CreateDate,
			&user.ImagePath,
			&user.Verified,
		)

		if err != nil {
			log.Println(err)
			return User{}, err
		}

		rows, err = conn.Query(context.Background(), "SELECT link_url, image_path from links WHERE user_id = $1", userId)

		if err != nil {
			log.Println(err)
			return User{}, err
		}

		for rows.Next() {
			lwp := new(LinkWithPath)

			err = rows.Scan(&lwp.Link, &lwp.Path)

			if err != nil {
				log.Println(err)
				return User{}, err
			}

			user.Links = append(user.Links, []string{lwp.Link, lwp.Path})
		}

		if rows.Err() != nil {
			log.Println(rows.Err())
			return User{}, rows.Err()
		}

		err = conn.QueryRow(context.Background(), "SELECT COUNT(*) from likes WHERE user_likes_id = $1", userId).Scan(
			&user.LikesCount,
		)

		if err != nil {
			log.Println(err)
			return User{}, err
		}

		if len(user.ImagePath) == 0 {
			user.ImagePath = constants.UserImages + constants.DefaultUserImage
		}
	}

	return user, nil
}

func AddUser(conn *pgx.Conn, username string, description string, password string, email string, imagePath string) (int, error) {
	var userId int

	err := conn.QueryRow(context.Background(),
		"INSERT INTO users (username, description, password, email, create_date, image_path, verified)"+
			"VALUES ($1, $2, $3, $4, $5, $6, false) RETURNING user_id", username, description, password, email, time.Now(), imagePath).Scan(&userId)

	if err != nil {
		return -1, err
	}

	return userId, nil
}

func CheckEmailExistInDB(conn *pgx.Conn, email string) (bool, error) {
	var emailDB string

	err := conn.QueryRow(context.Background(), "SELECT email from users WHERE email=$1", email).Scan(&emailDB)

	if err != nil {
		log.Println(err)
		return false, err
	}

	return emailDB == email, nil
}

func GetUserByEmailAndPassword(conn *pgx.Conn, email string, password string) (User, error) {
	var user User
	var userIdDB int
	var passwordDB string
	var err error

	{
		err = conn.QueryRow(context.Background(), "SELECT user_id, password from users WHERE email=$1", email).Scan(
			&userIdDB,
			&passwordDB,
		)

		if err != nil {
			return User{}, err
		}
	}

	{
		PasswordsCompare, err := CheckPassword(password, passwordDB)

		if err != nil || !PasswordsCompare {
			return User{}, ErrWrongPassword
		}
	}

	{
		user, err = GetUserViaId(conn, userIdDB)
		if err != nil {
			return User{}, err
		}
	}

	return user, nil
}

func UpdateUser(conn *pgx.Conn, user User) error {
	var userId int

	err := conn.QueryRow(context.Background(), "UPDATE users SET "+
		"username=$1, email=$2, password=$3, image_path=$4, verified=$5 WHERE user_id=$6 RETURNING user_id",
		user.Username, user.Email, user.Password, user.ImagePath, user.Verified, user.UserId).Scan(&userId)

	return err
}
