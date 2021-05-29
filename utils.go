package main

import (
	"golang.org/x/crypto/bcrypt"
	"net/url"
	"strings"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func createIconLinkPairs(links []string) [][]string {
	var data [][]string

	for _, link := range links {
		allUrl, err := url.Parse(link)

		if err != nil {
			panic(err)
		}

		host := allUrl.Host

		if strings.Count(host, ".") >= 2 {

			for strings.Count(host, ".") >= 2 {
				host = host[strings.Index(host, ".")+1:]
			}

		}

		file := linksImagesPairs[host]

		if len(file) == 0 {
			file = linksImagesPairs["other"]
		}

		file = linkPath + file
		data = append(data, []string{link, file})
	}

	return data
}
