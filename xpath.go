package main

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"net/http"
	"strings"
)

type Department struct {
	Name string          `json:"name"`
	Days map[string]*Day `json:"days"`
}

type Day struct {
	Name    string    `json:"name"`
	Lessons []*Lesson `json:"lessons"`
}

type Lesson struct {
	ID     string            `json:"id"`
	Time   string            `json:"time"`
	Type   string            `json:"type"`
	Week   string            `json:"week"`
	Name   string            `json:"name"`
	Tutors map[string]string `json:"tutors"`
	Groups map[string]string `json:"groups"`
	Room   string            `json:"room"`
	RoomID string            `json:"room_id"`
}

func parseByXpath(url string) {
	newTerms := make(map[string]string)
	resp, err := http.Get(url)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		// handle error
	}

	departmentNameNode := htmlquery.FindOne(doc, "/html/body/div[1]/div/div/div[3]/h1")
	departmentName := strings.TrimSpace(departmentNameNode.FirstChild.Data)
	newTerms["departmentName"] = departmentName

	days := htmlquery.Find(doc, "/html/body/div[1]/div/div/div[contains(@class,'list-group')]")
	dayNames := htmlquery.Find(doc, "/html/body/div[1]/div/div/h3[@class = 'lesson-wday']")

	for i, day := range days {
		if i < len(dayNames) {
			dayName := strings.TrimSpace(dayNames[i].FirstChild.Data)
			fmt.Println(dayName)
		}
		lessonsGroupItem := htmlquery.Find(day, "./div[@class = 'list-group-item']")

		for _, lessonGroupItem := range lessonsGroupItem {
			lessonTimeNode := htmlquery.FindOne(lessonGroupItem, "./div[@class = 'lesson-time']")
			lessonTime := strings.TrimSpace(lessonTimeNode.FirstChild.Data)
			lessonTime = strings.ReplaceAll(lessonTime, "—", "-")
			fmt.Println(lessonTime)

			lessons := htmlquery.Find(lessonGroupItem, "./div[@class = 'lesson-lessons']/div")
			// Process lessons
			for _, lesson := range lessons {
				// print 'data-id' attribute
				lessonID := htmlquery.SelectAttr(lesson, "data-id")
				fmt.Println(lessonID)
				lessonRoomNode := htmlquery.FindOne(lesson, "./div/a/text()")
				lessonRoom := ""
				lessonRoomId := ""
				if lessonRoomNode != nil {
					lessonRoom = strings.TrimSpace(lessonRoomNode.Data)
					lessonRoomId = htmlquery.SelectAttr(lessonRoomNode.Parent, "href")
					lessonRoomId = strings.ReplaceAll(lessonRoomId, "/rooms/", "")
					fmt.Println(lessonRoom, lessonRoomId)
				} else {
					lessonRoomNode = htmlquery.FindOne(lesson, "./div/span/text()")
					lessonRoom = strings.TrimSpace(lessonRoomNode.Data)
					lessonRoomId = ""
					fmt.Println(lessonRoom, lessonRoomId)
				}

			}
		}
	}
}
