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
	Dev         bool
	Status      int
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
			&user.Dev,
			&user.Status,
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
		"INSERT INTO users (username, description, password, email, create_date, image_path, verified, dev, status)"+
			"VALUES ($1, $2, $3, $4, $5, $6, false, false, 1) RETURNING user_id", username, description, password, email, time.Now(), imagePath).Scan(&userId)

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

func UpdateUserMainInfo(conn *pgx.Conn, user User) error {
	var userId int

	err := conn.QueryRow(context.Background(), "UPDATE users SET "+
		"username=$1, description=$2, email=$3, password=$4, image_path=$5 WHERE user_id=$6 RETURNING user_id",
		user.Username, user.Description, user.Email, user.Password, user.ImagePath, user.UserId).Scan(&userId)

	return err
}

func UpdateUserDevInfo(conn *pgx.Conn, user User) error {
	var userId int

	err := conn.QueryRow(context.Background(), "UPDATE users SET "+
		"verified=$1, dev=$2, status=$3 WHERE user_id=$4 RETURNING user_id",
		user.Verified, user.Dev, user.Status, user.UserId).Scan(&userId)

	return err
}

func UpdateUsersLinks(conn *pgx.Conn, userId int, links [][]string) error {
	var answer int
	err := conn.QueryRow(context.Background(), "DELETE FROM links WHERE user_id=$1 RETURNING 1", userId).Scan(&answer)

	if err != nil {
		if err.Error() != "no rows in result set" {
			return err
		}
	}

	var rows [][]interface{}

	for _, row := range links {
		rows = append(rows, []interface{}{row[0], row[1], userId})
	}

	_, err = conn.CopyFrom(context.Background(), pgx.Identifier{"links"},
		[]string{"link_url", "image_path", "user_id"}, pgx.CopyFromRows(rows))

	return err
}

func SetLike(conn *pgx.Conn, userLikedId, userLikesId int) error {
	err := conn.QueryRow(context.Background(),
		"INSERT INTO likes (user_liked_id, user_likes_id) VALUES ($1, $2)", userLikedId, userLikesId).Scan()

	if err.Error() != "no rows in result set" {
		return err
	}

	return nil
}

func DeleteLike(conn *pgx.Conn, userLikedId, userLikesId int) error {
	err := conn.QueryRow(context.Background(),
		"DELETE FROM likes WHERE user_liked_id=$1 AND user_likes_id=$2", userLikedId, userLikesId).Scan()

	if err.Error() != "no rows in result set" {
		return err
	}

	return nil
}
