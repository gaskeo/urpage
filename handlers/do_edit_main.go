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
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func CreateDoEditMain(conn *pgx.Conn, rds *redis.Client) {

	doEditMain := func(writer http.ResponseWriter, request *http.Request) {
		var (
			userId                                        int
			username, imageName, CSRFToken, CSRFTokenForm string
			jsonAnswer                                    []byte
			newImage                                      *os.File
			user                                          structs.User
			imageForm                                     multipart.File
			err                                           error
		)

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

			username = request.FormValue("username")
			imageForm, _, err = request.FormFile("image") // header with name
			// check format of file
			if err == nil {
				defer func() {
					err = imageForm.Close()
					if err != nil {
						jsonAnswer, _ = json.Marshal(structs.Answer{Err: "other-error"})
						return
					}
				}()

				imageName, err = utils.GenerateImageName()

				if err != nil {
					jsonAnswer, _ = json.Marshal(structs.Answer{Err: "other-error"})
					return
				}

				_, err = os.Create(constants.UserImages[1:] + imageName + ".jpeg")

				if err != nil {
					jsonAnswer, _ = json.Marshal(structs.Answer{Err: "other-error"})
					return
				}
				// todo check image type
				newImage, err = os.OpenFile(constants.UserImages[1:]+imageName+".jpeg", os.O_WRONLY, 0644)

				if err != nil {
					jsonAnswer, _ = json.Marshal(structs.Answer{Err: "other-error"})
					return
				}

				defer func() {
					err = newImage.Close()

					if err != nil {

						jsonAnswer, _ = json.Marshal(structs.Answer{Err: "other-error"})
						return
					}
				}()

				_, err = io.Copy(newImage, imageForm)

				if err != nil {

					jsonAnswer, _ = json.Marshal(structs.Answer{Err: "other-error"})
					return
				}
			}
		}

		{ // get user
			user, err = storage.GetUserViaId(conn, userId)

			if err != nil {
				http.Error(writer, "Ошибка с БД", http.StatusForbidden)
				return
			}
		}

		{ // set new data
			if len(imageName) > 0 {
				user.ImagePath = imageName + ".jpeg"
			} else {
				user.ImagePath = user.ImagePath[len(constants.UserImages):]
			}

			if len(username) > 0 {
				user.Username = username
			}

			err = storage.UpdateUser(conn, user)

			if err != nil {
				jsonAnswer, _ = json.Marshal(structs.Answer{Err: "other-error"})
				return
			}

			jsonAnswer, _ = json.Marshal(structs.Answer{Err: ""})
		}
	}

	http.HandleFunc("/do/edit_main", doEditMain)
}
