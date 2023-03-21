package parser

import (
	"encoding/json"
	"fmt"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type RoomId struct {
	Id int
}

type Room struct {
	Id           int
	Lab          bool
	Projector    bool
	Computer     bool
	Availability bool
	DOT          bool
	Temporary    bool
}

func ParseAllRooms() {
	newTerms := make(map[int][]string)
	rooms := ParseRoomsJson()
	counter := 0
	fmt.Println(len(rooms))
	for _, element := range rooms {
		roomData := ParseRoomByXpath("https://home.mephi.ru/rooms/" + strconv.Itoa(element))
		counter++
		if roomData == nil {
			newTerms[element] = []string{""}
		} else {
			newTerms[element] = roomData
		}
		if counter%300 == 0 {
			jsonData, err := json.MarshalIndent(newTerms, "", " ")
			if err != nil {
				fmt.Println(err)
			}
			_ = ioutil.WriteFile("rooms.json", jsonData, 0644)
			fmt.Printf("Parsed %d rooms\n", counter)
		}
	}
	jsonData, err := json.MarshalIndent(newTerms, "", " ")
	if err != nil {
		fmt.Println(err)
	}
	_ = ioutil.WriteFile("rooms.json", jsonData, 0644)
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

func CreateInserts() {
	data, err := ioutil.ReadFile("./ydb/sources/parsed/rooms.json")
	var parsedData map[int][]string
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(data), &parsedData)
	if err != nil {
		log.Fatal(err)
	}
	for key, value := range parsedData {
		// split value by " "
		//fmt.Printf("INSERT INTO rooms (id, lab, projector, computer, availability, dot, temporary) VALUES (%d, %t, %t, %t, %t, %t, %t);\n", key, value[0] == "lab", value[0] == "projector", value[0] == "computer", value[0] == "availability", value[0] == "dot", value[0] == "temporary")
	}
}
