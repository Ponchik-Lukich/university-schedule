package parser

import (
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"net/http"
	"strings"
)

func ParseRoomByXpath(url string) map[string][]string {
	//newTerms := make(map[string]interface{})
	var tags []string

	resp, err := http.Get(url)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		// handle error
	}
	roomId := strings.ReplaceAll(url, "https://home.mephi.ru/rooms/", "")
	tagsNode := htmlquery.Find(doc, "/html/body/div/div/div/div[3]/h1/span/i")
	for _, node := range tagsNode {
		tag := node.Attr[0].Val
		tags = append(tags, tag)
	}
	roomData := make(map[string][]string)
	roomData[roomId] = tags
	return roomData
}
