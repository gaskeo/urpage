package utils

import (
	"go-site/constants"
	"golang.org/x/crypto/bcrypt"
	"net/url"
	"strings"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPassword(password string, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil
}

func CreateIconLinkPairs(links []string) ([][]string, error) {
	var data [][]string

	for _, link := range links {
		allUrl, err := url.Parse(link)

		if err != nil {
			return [][]string{}, err
		}

		host := allUrl.Host

		if strings.Count(host, ".") >= 2 {

			for strings.Count(host, ".") >= 2 {
				host = host[strings.Index(host, ".")+1:]
			}

		}

		file := constants.LinksImagesPairs[host]

		if len(file) == 0 {
			file = constants.LinksImagesPairs["other"]
		}

		file = constants.LinkPath + file
		data = append(data, []string{link, file})
	}

	return data, nil
}
