package main

import (
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
	newTerms := make(map[string]map[string]interface{})
	resp, err := http.Get(url)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		// handle error
	}

	departmentId := strings.ReplaceAll(url, "https://home.potatohd.ru/departments/", "")
	departmentNameNode := htmlquery.FindOne(doc, "/html/body/div[1]/div/div/div[3]/h1")
	departmentName := strings.TrimSpace(departmentNameNode.FirstChild.Data)
	newTerms[departmentId] = make(map[string]interface{})
	newTerms[departmentId]["days"] = make(map[string]interface{})
	newTerms[departmentId][departmentId] = departmentName

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
			lessonTime = strings.ReplaceAll(lessonTime, "—", "-")
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
					tutorId := htmlquery.SelectAttr(tutor, "href")
					tutorId = strings.ReplaceAll(tutorId, "/tutors/", "")
					tutorsData[tutorId] = tutorName
					//fmt.Println(tutorName, tutorId)
				}

				// get lesson dates
				lessonDatesNode := htmlquery.FindOne(lesson, "./span[@class = 'lesson-dates']/text()")
				lessonDates := ""
				if lessonDatesNode != nil {
					lessonDates = strings.TrimSpace(lessonDatesNode.Data)
					lessonDates = strings.ReplaceAll(lessonDates, "—", "-")
					//fmt.Println(lessonDates)
				}
				additionalInfoNode := htmlquery.Find(lesson, "./span[@class = 'text-muted']/text()")
				additionalInfo := ""
				if len(additionalInfoNode) > 0 {
					for _, info := range additionalInfoNode {
						additionalInfo += strings.TrimSpace(info.Data)
					}
					//fmt.Println(additionalInfo)
				}
				// new lessonsData
				lessonData := LessonData{
					Time:   lessonTime,
					Type:   lessonType,
					Week:   lessonWeek,
					Name:   lessonName,
					Tutors: tutorsData,
					Groups: groupsData,
					Room:   lessonRoom,
					RoomID: lessonRoomId,
				}
				dayData[lessonID] = lessonData
			}
			newTerms[departmentId]["days"].(map[string]interface{})[dayName] = dayData
		}
	}
	// print all data beautifully
	//
	//fmt.Println(newTerms[departmentId]["days"].(map[string]interface{})["Понедельник"].(map[string]LessonData)["402882"].Time)
	//fmt.Println(newTerms[departmentId]["days"].(map[string]interface{})["Понедельник"].(map[string]LessonData)["1"])

}
