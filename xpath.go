package main

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/antchfx/xpath"
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

	root := htmlquery.CreateXPathNavigator(doc)
	path := "/html/body/div[1]/div/div/div[3]/h1/text()"
	expr := xpath.MustCompile(path)

	departmentNameRes := expr.Evaluate(root).(*xpath.NodeIterator)
	departmentNameRes.MoveNext()
	node := departmentNameRes.Current()
	departmentName := node.Value()
	departmentName = strings.ReplaceAll(departmentName, "\n", "")
	departmentName = strings.TrimSpace(departmentName)
	newTerms["departmentName"] = departmentName
	//newTerms["days"] = make(map[string]string)

	path = "/html/body/div[1]/div/div/div[contains(@class,'list-group')]"
	expr = xpath.MustCompile(path)
	days := expr.Evaluate(root).(*xpath.NodeIterator)
	path = "/html/body/div[1]/div/div/h3[@class = 'lesson-wday']/text()"
	expr = xpath.MustCompile(path)

	dayNames := expr.Evaluate(root).(*xpath.NodeIterator)

	for days.MoveNext() {

		path = "./div[@class = 'list-group-item']"
		expr = xpath.MustCompile(path)
		lessonsGroupItem := expr.Evaluate(days.Current()).(*xpath.NodeIterator)

		//dayData := make(map[string]string)

		for lessonsGroupItem.MoveNext() {

			path = "./div[@class = 'lesson-time']/text()"
			expr = xpath.MustCompile(path)
			lessonTimeRes := expr.Evaluate(lessonsGroupItem.Current()).(*xpath.NodeIterator)
			node := lessonTimeRes.Current()
			lessonTime := node.Value()
			lessonTime = strings.ReplaceAll(lessonTime, "\n", "")
			lessonTime = strings.TrimSpace(lessonTime)
			lessonTime = strings.ReplaceAll(lessonTime, "—", "-")
			fmt.Println("HERE1", lessonTimeRes.Current().Value())
			path = "./div[@class = 'lesson-lessons']/div"
			expr = xpath.MustCompile(path)
			lessonsRes := expr.Evaluate(lessonsGroupItem.Current()).(*xpath.NodeIterator)

		}
		if dayNames.MoveNext() {
			fmt.Println(dayNames.Current().Value())
		}

	}

}
