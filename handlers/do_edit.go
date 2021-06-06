package handlers

import (
	"fmt"
	"go-site/constants"
	"go-site/storage"
	"go-site/utils"
	"go-site/verify_utils"
	"io"
	"log"
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

	if request.Method != "POST" {
		return
	}

	{ // check user authed
		userId, err = verify_utils.CheckIfUserAuth(request)

		if err != nil {
			log.Println(err)
			http.Redirect(writer, request, "/", http.StatusSeeOther)
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
					log.Println(err)
				}
			}()

			imageName, err = utils.GenerateImageName()

			if err != nil {
				panic(err)
				return
			}

			_, err = os.Create(constants.UserImages[1:] + imageName + ".jpeg")

			if err != nil {
				log.Println(err, 1)
			}

			newImage, err := os.OpenFile(constants.UserImages[1:]+imageName+".jpeg", os.O_WRONLY, 0644)

			if err != nil {
				panic(err)
			}

			defer func() {
				err = newImage.Close()

				if err != nil {
					log.Println(err)
				}
			}()

			_, err = io.Copy(newImage, imageForm)

			if err != nil {
				log.Println(err, 2)
			}
		} else {
			log.Println(err)
		}
	}

	{ // compare form and authed user
		userIdIntForm, err := strconv.Atoi(userIdStr)

		if err != nil {
			log.Println(err, 13)
		}

		if userId != userIdIntForm {
			http.Redirect(writer, request, "/", http.StatusSeeOther)
			return
		}
	}

	{ // get user
		user, err = storage.GetUserViaId(userId)

		if err != nil {
			log.Println(err)
		}
	}

	{ // set new data
		if len(imageName) > 0 {
			fmt.Println(imageName)
			user.ImagePath = imageName + ".jpeg"
		} else {
			user.ImagePath = user.ImagePath[len(constants.UserImages):]
		}

		if len(username) > 0 {
			user.Username = username
		}

		err = storage.UpdateUser(user)

		if err != nil {
			panic(err)
			return
		}

	}
}
