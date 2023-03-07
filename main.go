package main

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
	"university-timetable/hash"
)

var websites = []string{
	"https://home.potatohd.ru/departments/2603786",
}

var prevContent = []string{
	"",
}

var prevXpath = []string{
	"",
}

var weekDays = []string{
	"Понедельник",
	"Вторник",
	"Среда",
	"Четверг",
	"Пятница",
	"Суббота",
	"Воскресенье",
}

var classes = []string{
	"text-nowrap",
	"lesson-square lesson-square-0",
	"lesson-square lesson-square-1",
	"lesson-square lesson-square-2",
}

var wg sync.WaitGroup

func main() {
	//connect()

	hash.GetHash()

	//parser.ParseByXpath("https://home.potatohd.ru/departments/2603786")
	//parser.ParseByXpathExam("https://home.potatohd.ru/departments/111056/exams")

	//for i, url := range websites {
	//	wg.Add(1)
	//	go func(i int, url string) {
	//		defer wg.Done()
	//		result := false
	//		data, _ := getWebsiteData(url)
	//		doc, _ := goquery.NewDocumentFromReader(strings.NewReader(data))
	//		body := doc.Find("*").First()
	//		removeJunk(body, 0)
	//		if prevXpath[i] != "" {
	//			body = getXpathData(body, prevXpath[i])
	//		}
	//		data = body.Text()
	//		// remove all extra spaces
	//		data, _ = parseData(data)
	//		hash, _ := getWebsiteHash(data)
	//		if hash != prevContent[i] {
	//			prevContent[i] = hash
	//			result = true
	//		}
	//
	//		if result {
	//			fmt.Printf("Website %s changed!\n", url)
	//		} else {
	//			fmt.Printf("Nothing on %s\n", url)
	//		}
	//
	//	}(i, url)
	//}
	//wg.Wait()
}

func parseData(data string) (string, error) {
	re := regexp.MustCompile(`\n\n+`)
	data = re.ReplaceAllString(data, "\n")
	lines := strings.Split(data, "\n")
	var index int
out:
	for i, line := range lines {
		for _, day := range weekDays {
			if strings.HasPrefix(line, day) {
				index = i
				break out
			}
		}
	}
	lines = lines[index:]
	lines = lines[:len(lines)-3]
	data = strings.Join(lines, "\n")
	fmt.Println(data)
	return data, nil

}

func getWebsiteHash(data string) (string, error) {
	hasher := sha1.New()
	_, err := hasher.Write([]byte(data))
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil)), nil
}

func getWebsiteData(url string) (string, error) {
	method := "GET"

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return "", err
	}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic("Error closing body. Here's why: " + err.Error())
		}
	}(res.Body)
	if res.Status != "200 OK" {
		return "", errors.New("status code is not 200")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func removeJunk(data *goquery.Selection, choice int) {
	class, _ := data.Attr("class")
	for _, c := range classes {
		if class == c || choice == 1 {
			if class == "text-nowrap" {
				link := data.Find("a[href*='/tutors/']")
				if link.Length() > 0 {
					linkText := link.AttrOr("href", "")
					splitLink := strings.Split(linkText, "/tutors/")
					if len(splitLink) > 1 {
						data.SetText(splitLink[1])
						return
					}
				}
			} else if choice == 1 {
				if class == "lesson lesson-practice" || class == "lesson lesson-lecture" || class == "lesson lesson-lab" {
					dataID, exists := data.Attr("data-id")
					if exists {
						data.SetText(dataID)
						return
					}
				}
			} else {
				data.SetText(class)
				return
			}
		}
	}
	data = data.RemoveAttr("class")
	data = data.RemoveAttr("id")
	data.Contents().Each(func(i int, s *goquery.Selection) {
		if s.Is("script") {
			s.Remove()
		} else if s.Is("style") {
			s.Remove()
		} else {
			removeJunk(s, choice)
		}
	})
}

func setText(data *goquery.Selection, text string) {
	if data.Length() > 0 {
		content := data.Contents().Text()
		if len(content) > 0 {
			data.SetHtml(content + text)
		} else {
			data.SetText(text)
		}
	}
}

func getXpathData(body *goquery.Selection, xpath string) *goquery.Selection {
	xpath = strings.ReplaceAll(xpath, "/html/body", "")

	return body.Find(xpath)
}
