package handlers

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4"
	"go-site/constants"
	"go-site/jwt_api"
	"go-site/session"
	"go-site/storage"
	"go-site/utils"
	"net/http"
	"strings"
)

func CreateDoEditLinks(conn *pgx.Conn, rdb *redis.Client) {

	doEditLinks := func(writer http.ResponseWriter, request *http.Request) {
		var (
			userId           int
			links, CSRFToken string
			jsonAnswer       []byte
			user             storage.User
			err              error
		)

		if request.Method != "POST" {
			return
		}

		defer func() { SendJson(writer, jsonAnswer) }()

		{ // check is session has CSRF
			_, CSRFToken, err = session.CheckSessionId(writer, request, rdb)

			if err != nil {
				jsonAnswer, _ = json.Marshal(Answer{Err: "no-csrf"})
				return
			}
		}

		{ // check user authed
			userId, err = jwt_api.CheckIfUserAuth(writer, request, rdb)

			if err != nil {
				http.Error(writer, "no jwt", http.StatusForbidden)
				return
			}
		}

		{ // work with form and check CSRF
			CSRFTokenForm := request.FormValue("csrf")

			if CSRFToken != CSRFTokenForm {
				jsonAnswer, _ = json.Marshal(Answer{Err: "no-csrf"})
				return
			}

			links = request.FormValue("links")
		}

		{ // get user
			user, err = storage.GetUserViaId(conn, userId)

			if err != nil {
				http.Error(writer, "error getting user", http.StatusInternalServerError)
				return
			}
		}

		{ // set new data
			user.ImagePath = user.ImagePath[len(constants.UserImages):]

			linksLst := strings.Split(links, " ")

			user.Links, err = utils.CreateIconLinkPairs(linksLst)

			if err != nil {
				http.Error(writer, "error generating links", http.StatusInternalServerError)
				return
			}

			err = storage.UpdateUser(conn, user)

			if err != nil {
				http.Error(writer, "error updating user", http.StatusInternalServerError)
				return
			}
			jsonAnswer, _ = json.Marshal(Answer{Err: ""})

		}
	}

	http.HandleFunc("/do/edit_links", doEditLinks)
}
