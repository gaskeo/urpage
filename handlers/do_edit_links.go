package handlers

import (
	"fmt"
	"go-site/constants"
	"go-site/storage"
	"go-site/utils"
	"go-site/verify_utils"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func DoEditLinks(writer http.ResponseWriter, request *http.Request) {
	var user storage.User
	var userId int
	var userIdStr, links string
	var err error

	if request.Method != "POST" {
		return
	}

	{ // check user authed
		userId, err = verify_utils.CheckIfUserAuth(writer, request)

		if err != nil {
			log.Println(err)
			http.Redirect(writer, request, "/", http.StatusSeeOther)
			return
		}
	}

	{ // work with form
		userIdStr = request.FormValue("id")
		links = request.FormValue("links")
		fmt.Println(userIdStr, links, userId)
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
		user.ImagePath = user.ImagePath[len(constants.UserImages):]

		linksLst := strings.Split(links, " ")

		user.Links, err = utils.CreateIconLinkPairs(linksLst)

		if err != nil {
			panic(err)
			return
		}

		err = storage.UpdateUser(user)

		if err != nil {
			panic(err)
			return
		}
	}
}
