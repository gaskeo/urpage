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
	"io"
	"net/http"
	"os"
	"strings"
)

func CreateDoEditMain(conn *pgx.Conn, rdb *redis.Client) {

	doEditMain := func(writer http.ResponseWriter, request *http.Request) {
		var (
			userId                                    int
			username, imageName, CSRFToken, imageType string
			jsonAnswer                                []byte
			user                                      storage.User
			err                                       error
		)

		if request.Method != "POST" {
			return
		}

		defer func() { SendJson(writer, jsonAnswer) }()

		{ // CSRF check
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

		{ // work with form
			CSRFTokenForm := request.FormValue("csrf")

			if CSRFToken != CSRFTokenForm {
				jsonAnswer, _ = json.Marshal(Answer{Err: "no-csrf"})
				return
			}

			username = request.FormValue("username")
			imageForm, header, err := request.FormFile("image") // header with name
			// check format of file
			if err == nil {
				defer func() {
					err = imageForm.Close()
					if err != nil {
						http.Error(writer, "error generating image", http.StatusInternalServerError)
						return
					}
				}()

				imageName, err = utils.GenerateImageName()

				if err != nil {
					jsonAnswer, _ = json.Marshal(Answer{Err: "other-error"})
					return
				}

				fileSplit := strings.Split(header.Filename, ".")
				imageType = fileSplit[len(fileSplit)-1]

				switch imageType {
				case "jpg", "jpeg", "png":
					_, err = os.Create(constants.UserImages[1:] + imageName + "." + imageType)

					if err != nil {
						jsonAnswer, _ = json.Marshal(Answer{Err: "other-error"})
						return
					}
				default:
					jsonAnswer, _ = json.Marshal(Answer{Err: "bad-image-error"})
					return
				}

				newImage, err := os.OpenFile(constants.UserImages[1:]+imageName+"."+imageType, os.O_WRONLY, 0644)

				if err != nil {
					http.Error(writer, "error creating file", http.StatusForbidden)
					return
				}

				defer func() {
					err = newImage.Close()

					if err != nil {

						http.Error(writer, "error saving file", http.StatusForbidden)
						return
					}
				}()

				_, err = io.Copy(newImage, imageForm)

				if err != nil {

					http.Error(writer, "error getting file", http.StatusForbidden)
					return
				}
			}
		}

		{ // get user
			user, err = storage.GetUserViaId(conn, userId)

			if err != nil {
				http.Error(writer, "error getting user", http.StatusInternalServerError)
				return
			}
		}

		{ // delete exist image
			if user.ImagePath != constants.UserImages+"default.jpeg" && len(imageName) > 0 {
				err = os.Remove(user.ImagePath[1:])

				if err != nil {
					http.Error(writer, "error deleting old image", http.StatusInternalServerError)
				}
			}
		}

		{ // set new data
			if len(imageName) > 0 {
				user.ImagePath = imageName + "." + imageType
			}

			if len(username) > 0 {
				user.Username = username
			}

			err = storage.UpdateUser(conn, user)

			if err != nil {
				http.Error(writer, "error updating user", http.StatusInternalServerError)
				return
			}

			jsonAnswer, _ = json.Marshal(Answer{Err: ""})
		}
	}

	http.HandleFunc("/do/edit_main", doEditMain)
}
