package parser

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"net/http"
	"regexp"
	"strings"
)

type DayData struct {
	Name string
	Data map[string]LessonData
}

type LessonData struct {
	TimeFrom string            `json:"time_from"`
	TimeTo   string            `json:"time_to"`
	Type     string            `json:"type"`
	Week     string            `json:"week"`
	Name     string            `json:"name"`
	Tutors   map[string]string `json:"tutors"`
	Groups   map[string]string `json:"groups"`
	Room     string            `json:"room"`
	RoomID   string            `json:"room_id"`
	DateFrom string            `json:"date_from"`
	DateTo   string            `json:"date_to"`
	Dates    string            `json:"dates"`
	//Addition string            `json:"additional_info"`
}

func ParseByXpath(url string) {
	newTerms := make(map[string]map[string]interface{})

	daysHash := make(map[string]string)
	lessonsHash := make(map[string]string)

	resp, err := http.Get(url)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		// handle error
	}
	semester := htmlquery.FindOne(doc, "//*[@id=\"page-content-wrapper\"]/div/div[3]/ol/li[1]/a")
	semesterUrl := htmlquery.SelectAttr(semester, "href")
	semesterId := strings.ReplaceAll(semesterUrl, "/study_groups?organization_id=1&term_id=", "")
	departmentId := strings.ReplaceAll(url, "https://home.potatohd.ru/departments/", "")
	departmentNameNode := htmlquery.FindOne(doc, "/html/body/div[1]/div/div/div[3]/h1")
	departmentName := strings.TrimSpace(departmentNameNode.FirstChild.Data)
	newTerms[semesterId] = make(map[string]interface{})
	newTerms[semesterId]["days"] = make(map[string]interface{})
	//newTerms[departmentId][departmentName] = departmentName
	newTerms[semesterId][departmentId] = departmentName

	days := htmlquery.Find(doc, "/html/body/div[1]/div/div/div[contains(@class,'list-group')]")
	dayNames := htmlquery.Find(doc, "/html/body/div[1]/div/div/h3[@class = 'lesson-wday']")

	for i, day := range days {
		dayData := make(map[string]LessonData)
		dayName := strings.TrimSpace(dayNames[i].FirstChild.Data)
		//fmt.Println(dayName)
		lessonsGroupItem := htmlquery.Find(day, "./div[@class = 'list-group-item']")

		for _, lessonGroupItem := range lessonsGroupItem {

			lessonTimeNode := htmlquery.FindOne(lessonGroupItem, "./div[@class = 'lesson-time']")
			lessonTime := strings.TrimSpace(lessonTimeNode.FirstChild.Data)
			lessonTime = strings.ReplaceAll(lessonTime, " ", "")
			lessonTime = strings.ReplaceAll(lessonTime, " ", "")
			lessonTime = strings.ReplaceAll(lessonTime, "—", "-")
			lessonTimeFrom := strings.Split(lessonTime, "-")[0]
			lessonTimeTo := strings.Split(lessonTime, "-")[1]
			//fmt.Println(lessonTime)

			lessons := htmlquery.Find(lessonGroupItem, "./div[@class = 'lesson-lessons']/div")
			// Process lessons
			for _, lesson := range lessons {
				// print 'data-id' attribute
				lessonID := htmlquery.SelectAttr(lesson, "data-id")
				//fmt.Println(lessonID)
				lessonRoomNode := htmlquery.FindOne(lesson, "./div/a/text()")
				lessonRoom := ""
				lessonRoomId := ""
				if lessonRoomNode != nil {
					lessonRoom = strings.TrimSpace(lessonRoomNode.Data)
					lessonRoomId = htmlquery.SelectAttr(lessonRoomNode.Parent, "href")
					lessonRoomId = strings.ReplaceAll(lessonRoomId, "/rooms/", "")
					//fmt.Println(lessonRoom, lessonRoomId)
				} else {
					lessonRoomNode = htmlquery.FindOne(lesson, "./div/span/text()")
					lessonRoom = strings.TrimSpace(lessonRoomNode.Data)
					lessonRoomId = ""
					//fmt.Println(lessonRoom, lessonRoomId)
				}
				lessonWeekNode := htmlquery.FindOne(lesson, "./span[contains(@class, 'lesson-square')]")
				lessonWeek := htmlquery.SelectAttr(lessonWeekNode, "class")
				lessonWeek = strings.ReplaceAll(lessonWeek, "lesson-square lesson-square-", "")

				//fmt.Println(lessonWeek)
				lessonTypeNode := htmlquery.FindOne(lesson, "./div[contains(@class, 'label-lesson')]/text()")
				lessonType := ""
				if lessonTypeNode != nil {
					lessonType = strings.TrimSpace(lessonTypeNode.Data)
					//fmt.Println(lessonType)
				}
				lessonNameNode := htmlquery.Find(lesson, "./text()")
				lessonName := ""
				for _, node := range lessonNameNode {
					lessonName += strings.TrimSpace(node.Data)
				}
				re := regexp.MustCompile(`\s+`)
				lessonName = re.ReplaceAllString(lessonName, " ")
				lessonName = strings.TrimRight(lessonName, " ,")

				//fmt.Println(lessonName)

				// get groups
				groups := htmlquery.Find(lesson, "./a")
				groupsData := make(map[string]string)
				for _, group := range groups {
					groupName := strings.TrimSpace(group.FirstChild.Data)
					groupId := htmlquery.SelectAttr(group, "href")
					groupId = strings.ReplaceAll(groupId, "/study_groups/", "")
					groupId = strings.ReplaceAll(groupId, "/schedule", "")
					groupsData[groupId] = groupName
					//fmt.Println(groupName, groupId)
				}

				tutors := htmlquery.Find(lesson, "./span/a")
				tutorsData := make(map[string]string)
				for _, tutor := range tutors {
					tutorName := strings.TrimSpace(tutor.FirstChild.Data)
					tutorName = strings.ReplaceAll(tutorName, " ", " ")
					tutorId := htmlquery.SelectAttr(tutor, "href")
					tutorId = strings.ReplaceAll(tutorId, "/tutors/", "")
					tutorsData[tutorId] = tutorName
					//fmt.Println(tutorName, tutorId)
				}

				// get lesson dates
				lessonDatesNode := htmlquery.FindOne(lesson, "./span[@class = 'lesson-dates']/text()")
				lessonDate := ""
				lessonDatesFrom := ""
				lessonDatesTo := ""
				if lessonDatesNode != nil {
					lessonDates := strings.TrimSpace(lessonDatesNode.Data)
					lessonDates = strings.ReplaceAll(lessonDates, " ", "")
					lessonDates = strings.ReplaceAll(lessonDates, " ", "")
					lessonDates = strings.ReplaceAll(lessonDates, "—", "-")
					if strings.ContainsAny(lessonDates, "-") {
						lessonDates = strings.ReplaceAll(lessonDates, "(", "")
						lessonDates = strings.ReplaceAll(lessonDates, ")", "")
						lessonDates = strings.ReplaceAll(lessonDates, ",", "-")
						lessonDatesFrom = strings.Split(lessonDates, "-")[0]
						lessonDatesTo = strings.Split(lessonDates, "-")[1]
					} else {
						lessonDates = strings.ReplaceAll(lessonDates, "(", "")
						lessonDates = strings.ReplaceAll(lessonDates, ")", "")
						lessonDate = lessonDates
					}
					//fmt.Println(lessonDates)
				}
				lessonData := LessonData{
					TimeFrom: lessonTimeFrom,
					TimeTo:   lessonTimeTo,
					Type:     lessonType,
					Week:     lessonWeek,
					Name:     lessonName,
					Tutors:   tutorsData,
					Groups:   groupsData,
					Room:     lessonRoom,
					RoomID:   lessonRoomId,
					Dates:    lessonDate,
					DateFrom: lessonDatesFrom,
					DateTo:   lessonDatesTo,
					//Addition: additionalInfo,
				}
				// hash lessonData
				jsonLessonData, err := json.MarshalIndent(lessonData, "", "  ")
				if err != nil {
					fmt.Println(err)
				}
				//fmt.Println(string(jsonLessonData))
				// hash jsonData
				hash := sha256.Sum256([]byte(jsonLessonData))
				hashString := hex.EncodeToString(hash[:])
				// add to lessonsHash
				lessonsHash[lessonID] = hashString
				dayData[lessonID] = lessonData
			}
			newTerms[semesterId]["days"].(map[string]interface{})[dayName] = dayData
			// hash dayData
			jsondayData, err := json.MarshalIndent(dayData, "", " ")
			if err != nil {
				fmt.Println(err)
			}
			// hash jsonData
			hash := sha256.Sum256([]byte(jsondayData))
			hashString := hex.EncodeToString(hash[:])
			// add to daysHash
			daysHash[dayName] = hashString
		}
	}
	jsonData, err := json.MarshalIndent(newTerms[semesterId], "", " ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(jsonData))
	hash := sha256.Sum256([]byte(jsonData))
	hashString := hex.EncodeToString(hash[:])
	fmt.Println(hashString)
	fmt.Println(lessonsHash)
	fmt.Println(daysHash)
	//println(string(jsonData))
	//fmt.Println(newTerms[departmentId]["days"].(map[string]interface{})["Понедельник"].(map[string]LessonData)["402882"].Time)
	//fmt.Println(newTerms[departmentId]["days"].(map[string]interface{})["Понедельник"].(map[string]LessonData)["1"])
	//put data to json file
	// fill daysHash with
}
