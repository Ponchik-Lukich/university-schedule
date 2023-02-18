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

func zip(a, b []string) []string {
	var result []string
	for i := 0; i < len(a) && i < len(b); i++ {
		result = append(result, a[i]+" "+b[i])
	}
	return result
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

	departmentName := expr.Evaluate(root).(*xpath.NodeIterator)
	departmentName.MoveNext()
	node := departmentName.Current()
	dep := node.Value()
	dep = strings.ReplaceAll(dep, "\n", "")
	dep = strings.TrimSpace(dep)
	newTerms["departmentName"] = dep
	//newTerms["days"] = make(map[string]string)

	path = "/html/body/div[1]/div/div/div[contains(@class,'list-group')]"

	expr = xpath.MustCompile(path)

	days := expr.Evaluate(root).(*xpath.NodeIterator)

	path = "/html/body/div[1]/div/div/h3[@class = 'lesson-wday']/text()"

	expr = xpath.MustCompile(path)

	dayNames := expr.Evaluate(root).(*xpath.NodeIterator)

	var daysSlice, dayNamesSlice []string
	for days.MoveNext() {
		node := days.Current()
		day := node.Value()
		fmt.Println(day)
		daysSlice = append(daysSlice, node)
	}
	for dayNames.MoveNext() {
		node := dayNames.Current()
		dayName := node.Value()
		dayName = strings.ReplaceAll(dayName, "\n", "")
		dayName = strings.TrimSpace(dayName)
		dayNamesSlice = append(dayNamesSlice, node.Value())
	}

	zipped := zip(daysSlice, dayNamesSlice)

	// iterate over the zipped result
	for _, pair := range zipped {
		fmt.Println(pair)
	}

	// Iterate over the selected nodes
}
