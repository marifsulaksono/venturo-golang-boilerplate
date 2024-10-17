package helpers

import (
	"fmt"
	"net/url"
)

func LinksGenerator(base string, pageCount int, params map[string]interface{}) []string {
	links := []string{}

	for page := range pageCount {
		link := fmt.Sprintf("%s?", base)
		for param, value := range params {
			link += fmt.Sprintf("%s=%s&", url.QueryEscape(param), url.QueryEscape(fmt.Sprint(value)))
		}
		link += fmt.Sprintf("page=%d", page)

		links = append(links, link)
	}
	return links
}
