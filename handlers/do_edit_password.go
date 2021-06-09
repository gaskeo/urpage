package handlers

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4"
	"go-site/constants"
	"go-site/jwt"
	"go-site/session"
	"go-site/storage"
	"go-site/structs"
	"net/http"
)

func CreateDoEditPassword(conn *pgx.Conn, rds *redis.Client) {
	doEditPassword := func(writer http.ResponseWriter, request *http.Request) {
		var userId int
		var oldPassword, newPassword, CSRFToken, CSRFTokenForm string
		var correct bool

		var jsonAnswer []byte

		var user structs.User

		var err error

		if request.Method != "POST" {
			return
		}

		defer func() { SendJson(writer, jsonAnswer) }()

		{ // CSRF check
			_, CSRFToken, err = session.CheckSessionId(writer, request, rds)

			if err != nil {
				jsonAnswer, _ = json.Marshal(structs.Answer{Err: "no-csrf"})
				return
			}
		}

		{ // check user authed
			userId, err = jwt.CheckIfUserAuth(writer, request, rds)

			if err != nil {
				http.Error(writer, "У вас нет доступа", http.StatusForbidden)
				return
			}
		}

		{ // work with form
			CSRFTokenForm = request.FormValue("csrf")

			if CSRFToken != CSRFTokenForm {
				jsonAnswer, _ = json.Marshal(structs.Answer{Err: "no-csrf"})
				return
			}

			oldPassword = request.FormValue("old")
			newPassword = request.FormValue("new")
		}

		{ // get user
			user, err = storage.GetUserViaId(conn, userId)

			if err != nil {
				http.Error(writer, "Ошибка с БД", http.StatusForbidden)
				return
			}
		}

		{ // check old password
			correct, err = storage.CheckPassword(oldPassword, user.Password)
			if err != nil || !correct {
				jsonAnswer, _ = json.Marshal(structs.Answer{Err: "wrong-password"})
				return
			}
		}

		{ // set new data
			user.ImagePath = user.ImagePath[len(constants.UserImages):]

			user.Password, err = storage.HashPassword(newPassword)

			if err != nil {
				jsonAnswer, _ = json.Marshal(structs.Answer{Err: "other-error"})
				return
			}

			err = storage.UpdateUser(conn, user)
			if err != nil {
				jsonAnswer, _ = json.Marshal(structs.Answer{Err: "other-error"})
				return
			}
			jsonAnswer, _ = json.Marshal(structs.Answer{Err: ""})
		}
	}

	http.HandleFunc("/do/edit_password", doEditPassword)
}
