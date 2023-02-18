package main

import (
	"fmt"
	"log"
	"net/http"
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
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Parse HTML document
	//doc, err := htmlquery.Parse(resp.Body)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//var days []string
	//// Find the element using XPath
	//for i := 1; i < 6; i++ {
	//	day, err := htmlquery.Query(doc, "//*[@id=\"page-content-wrapper\"]/div/h3["+strconv.Itoa(i)+"]")
	//	if err != nil {
	//		continue
	//	} else {
	//		dayText := strings.TrimSpace(htmlquery.InnerText(day))
	//		if dayText != "" {
	//			days = append(days, dayText)
	//		}
	//	}
	//}

	for _, day := range days {
		fmt.Println(day)
	}
}
