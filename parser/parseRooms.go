package parser

import (
	"encoding/json"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"net/http"
)

type RoomId struct {
	Id int
}

func ParseRoomByXpath(url string) []string {
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
	tagsNode := htmlquery.Find(doc, "/html/body/div/div/div/div[3]/h1/span/i")
	for _, node := range tagsNode {
		tag := node.Attr[0].Val
		tags = append(tags, tag)
	}
	return tags
}

func ParseRoomsJson() []int {
	data, err := ioutil.ReadFile("./ydb/sources/parsed/mephi_public_rooms.json")
	var roomIds []int
	var parsedData []RoomId
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(data), &parsedData)
	if err != nil {
		log.Fatal(err)
	}
	for _, room := range parsedData {
		roomIds = append(roomIds, room.Id)
	}
	return roomIds
}
