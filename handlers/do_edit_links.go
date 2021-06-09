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
	"go-site/utils"
	"net/http"
	"strings"
)

func CreateDoEditLinks(conn *pgx.Conn, rds *redis.Client) {

	doEditLinks := func(writer http.ResponseWriter, request *http.Request) {
		var (
			userId                          int
			links, CSRFToken, CSRFTokenForm string
			linksLst                        []string
			jsonAnswer                      []byte
			user                            structs.User
			err                             error
		)

		if request.Method != "POST" {
			return
		}

		defer func() { SendJson(writer, jsonAnswer) }()

		{ // check is session has CSRF
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

		{ // work with form and check CSRF
			CSRFTokenForm = request.FormValue("csrf")

			if CSRFToken != CSRFTokenForm {
				jsonAnswer, _ = json.Marshal(structs.Answer{Err: "no-csrf"})
				return
			}

			links = request.FormValue("links")
		}

		{ // get user
			user, err = storage.GetUserViaId(conn, userId)

			if err != nil {
				http.Error(writer, "Ошибка с БД", http.StatusForbidden)
				return
			}
		}

		{ // set new data
			user.ImagePath = user.ImagePath[len(constants.UserImages):]

			linksLst = strings.Split(links, " ")

			user.Links, err = utils.CreateIconLinkPairs(linksLst)

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

	http.HandleFunc("/do/edit_links", doEditLinks)
}
