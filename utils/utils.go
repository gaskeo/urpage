package utils

import (
	"go-site/constants"
	"net/url"
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
