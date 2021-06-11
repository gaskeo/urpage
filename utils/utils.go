package utils

import (
	"crypto/rand"
	"go-site/constants"
	"math/big"
	"net/url"
	"strconv"
	"strings"
)

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

func CreateDBLinksFromPairs(LinkPairs [][]string) string {
	var DBLinks []string

	for _, linkPair := range LinkPairs {
		DBLinks = append(DBLinks, linkPair[0])
	}

	return strings.Join(DBLinks, " ")
}

func GenerateImageName() (string, error) {
	id, err := rand.Int(rand.Reader, big.NewInt(1000000000))

	if err != nil {
		return "", err
	}

	idStr := strconv.FormatInt(id.Int64(), 10)

	return idStr, nil
}
