package parser

import (
	"encoding/json"
	"fmt"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
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
	counter := 0
	for key, value := range parsedData {
		var lab, projector, computer, availability, dot, temporary bool
		for _, str := range value {
			if strings.Contains(str, "Лаборатория") {
				lab = true
			}
			if strings.Contains(str, "Аудитория оборудована проектором") {
				projector = true
			}
			if strings.Contains(str, "Компьютерный класс") {
				computer = true
			}
			if strings.Contains(str, "Аудитория общего фонда") {
				availability = true
				counter++
			}
			if strings.Contains(str, "Проводится с использованием дистанционных образовательных технологий") {
				dot = true
			}
			if strings.Contains(str, "Выездная аудитория") {
				temporary = true
			}
		}
		f, err := os.OpenFile("./ydb/sources/parsed/insert_rooms.sql", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			panic(err)
		}

		defer f.Close()

		if _, err = f.WriteString(fmt.Sprintf("UPDATE rooms SET (lab, projector, computer, availability, dot, temporary)="+
			"(%t, %t, %t, %t, %t, %t) WHERE id=%d;\n", lab, projector, computer, availability, dot, temporary, key)); err != nil {
			panic(err)
		}
		//_ = ioutil.WriteFile("./ydb/sources/parsed/insert_rooms.sql", []byte(fmt.Sprintf("UPDATE rooms SET (lab, projector, computer, availability, dot, temporary)="+
		//	"(%t, %t, %t, %t, %t, %t) WHERE id=%d;\n", lab, projector, computer, availability, dot, temporary, key)), 0644)
	}
	println(counter)
}
