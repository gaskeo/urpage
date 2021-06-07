package handlers

import (
	"encoding/json"
	"go-site/constants"
	"go-site/storage"
	"go-site/utils"
	"go-site/verify_utils"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

func DoEditMain(writer http.ResponseWriter, request *http.Request) {
	var userId int
	var userIdStr, username, imageName string
	var user storage.User
	var imageForm multipart.File
	var err error

	var jsonAnswer []byte

	if request.Method != "POST" {
		return
	}

	defer func() {
		if len(jsonAnswer) > 0 {
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write(jsonAnswer)
		}
	}()

	{ // check user authed
		userId, err = verify_utils.CheckIfUserAuth(writer, request)

		if err != nil {
			http.Error(writer, "У вас нет доступа", http.StatusForbidden)
			return
		}
	}

	{ // work with form
		userIdStr = request.FormValue("id")
		username = request.FormValue("username")
		imageForm, _, err = request.FormFile("image") // header with name
		// check format of file
		if err == nil {
			defer func() {
				err := imageForm.Close()
				if err != nil {
					jsonAnswer, _ = json.Marshal(Answer{Err: "other-error"})
					return
				}
			}()

			imageName, err = utils.GenerateImageName()

			if err != nil {
				jsonAnswer, _ = json.Marshal(Answer{Err: "other-error"})
				return
			}

			_, err = os.Create(constants.UserImages[1:] + imageName + ".jpeg")

			if err != nil {
				jsonAnswer, _ = json.Marshal(Answer{Err: "other-error"})
				return
			}

			newImage, err := os.OpenFile(constants.UserImages[1:]+imageName+".jpeg", os.O_WRONLY, 0644)

			if err != nil {
				jsonAnswer, _ = json.Marshal(Answer{Err: "other-error"})
				return
			}

			defer func() {
				err = newImage.Close()

				if err != nil {

					jsonAnswer, _ = json.Marshal(Answer{Err: "other-error"})
					return
				}
			}()

			_, err = io.Copy(newImage, imageForm)

			if err != nil {

				jsonAnswer, _ = json.Marshal(Answer{Err: "other-error"})
				return
			}
		}
	}

	{ // compare form and authed user
		userIdIntForm, err := strconv.Atoi(userIdStr)

		if err != nil {
			jsonAnswer, _ = json.Marshal(Answer{Err: "other-error"})
		}

		if userId != userIdIntForm {
			http.Error(writer, "У вас нет доступа", http.StatusForbidden)
			return
		}
	}

	{ // get user
		user, err = storage.GetUserViaId(userId)

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

		err = storage.UpdateUser(user)

		if err != nil {
			jsonAnswer, _ = json.Marshal(Answer{Err: "other-error"})
			return
		}

		jsonAnswer, _ = json.Marshal(Answer{Err: ""})
	}
}
