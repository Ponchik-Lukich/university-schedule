package parser

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"net/http"
	"regexp"
	"strings"
)

func convertDay(day string) string {
	switch day {
	case "Пн":
		day = "1"
	case "Вт":
		day = "2"
	case "Ср":
		day = "3"
	case "Чт":
		day = "4"
	case "Пт":
		day = "5"
	case "Сб":
		day = "6"
	case "Вс":
		day = "7"
	}
	return day
}

func convertMonth(month string) string {
	switch month {
	case "янв":
		month = "01"
	case "февр":
		month = "02"
	case "марта":
		month = "03"
	case "апр":
		month = "04"
	case "мая":
		month = "05"
	case "июня":
		month = "06"
	case "июля":
		month = "07"
	case "авг":
		month = "08"
	case "сент":
		month = "09"
	case "окт":
		month = "10"
	case "нояб":
		month = "11"
	case "дек":
		month = "12"
	}
	return month
}

func ParseByXpathExam(url string) {
	//newTerms := make(map[string]map[string]interface{})
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
	fmt.Println(semesterId)

	departmentId := strings.ReplaceAll(url, "https://home.potatohd.ru/departments/", "")
	departmentId = strings.ReplaceAll(departmentId, "/exams", "")
	departmentNameNode := htmlquery.FindOne(doc, "/html/body/div[1]/div/div/div[3]/h1")
	departmentName := strings.TrimSpace(departmentNameNode.FirstChild.Data)
	fmt.Println(departmentName, departmentId)

	lessomGroupItems := htmlquery.Find(doc, "/html/body/div/div/div/div/div[@class = 'list-group-item']/div")
	for _, lesson := range lessomGroupItems {

		lessonDateTimeNode := htmlquery.FindOne(lesson, "./div[@class = 'lesson-date']")
		lessonDateTime := strings.TrimSpace(lessonDateTimeNode.FirstChild.Data)
		// cut before ","
		weekDay := convertDay(strings.Split(lessonDateTime, ",")[0])
		fmt.Println(weekDay)
		// cut after "," before "."
		lessonDate := strings.Split(lessonDateTime, ",")[1]
		lessonDate = strings.Split(lessonDate, ".")[0]
		lessonDate = strings.TrimSpace(lessonDate)
		re := regexp.MustCompile(`(\d{2}\s+)([а-яА-Я]+)`)
		lessonDateParts := re.FindStringSubmatch(lessonDate)
		// convert month to number
		month := convertMonth(lessonDateParts[2])
		lessonDate = strings.TrimSpace(lessonDateParts[1]) + "." + month
		fmt.Println(lessonDate)

		lessonRoomNode := htmlquery.FindOne(lesson, "./div/a/text()")
		lessonRoom := ""
		lessonRoomId := ""
		if lessonRoomNode != nil {
			lessonRoom = strings.TrimSpace(lessonRoomNode.Data)
			lessonRoomId = htmlquery.SelectAttr(lessonRoomNode.Parent, "href")
			lessonRoomId = strings.ReplaceAll(lessonRoomId, "/rooms/", "")
			lessonRoomId = strings.ReplaceAll(lessonRoomId, "/exams", "")
		} else {
			lessonRoomNode = htmlquery.FindOne(lesson, "./div/span/text()")
			lessonRoom = strings.TrimSpace(lessonRoomNode.Data)
			lessonRoomId = ""
		}
		fmt.Println(lessonRoom, lessonRoomId)
		lessonTypeNode := htmlquery.FindOne(lesson, "./div[contains(@class, 'label-lesson')]/text()")
		lessonType := ""
		if lessonTypeNode != nil {
			lessonType = strings.TrimSpace(lessonTypeNode.Data)
			fmt.Println(lessonType)
		}
		lessonNameNode := htmlquery.Find(lesson, "./text()")
		lessonName := ""
		for _, node := range lessonNameNode {
			lessonName += strings.TrimSpace(node.Data)
		}
		re = regexp.MustCompile(`\s+`)
		lessonName = re.ReplaceAllString(lessonName, " ")
		lessonName = strings.TrimRight(lessonName, " ,")
		fmt.Println(lessonName)
		groups := htmlquery.Find(lesson, "./a")
		groupsData := make(map[string]string)
		for _, group := range groups {
			groupName := strings.TrimSpace(group.FirstChild.Data)
			groupId := htmlquery.SelectAttr(group, "href")
			groupId = strings.ReplaceAll(groupId, "/study_groups/", "")
			groupId = strings.ReplaceAll(groupId, "/schedule", "")
			groupId = strings.ReplaceAll(groupId, "/exams", "")
			groupsData[groupId] = groupName
			fmt.Println(groupName, groupId)
		}

		tutors := htmlquery.Find(lesson, "./span/a")
		tutorsData := make(map[string]string)
		for _, tutor := range tutors {
			tutorName := strings.TrimSpace(tutor.FirstChild.Data)
			tutorId := htmlquery.SelectAttr(tutor, "href")
			tutorId = strings.ReplaceAll(tutorId, "/tutors/", "")
			tutorId = strings.ReplaceAll(tutorId, "/exams", "")
			tutorsData[tutorId] = tutorName
			fmt.Println(tutorName, tutorId)
		}

		//lessonData := LessonData{
		//	TimeFrom: lessonTimeFrom,
		//	TimeTo:   lessonTimeTo,
		//	Type:     lessonType,
		//	Week:     "",
		//	Name:     lessonName,
		//	Tutors:   tutorsData,
		//	Groups:   groupsData,
		//	Room:     lessonRoom,
		//	RoomID:   lessonRoomId,
		//	Dates:    lessonDate,
		//	DateFrom: "",
		//	DateTo:   "",
		//	Addition: "",
		//}
	}

}
