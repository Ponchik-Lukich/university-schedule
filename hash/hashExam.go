package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
	"university-timetable/parser"
)

type DepartmentExams struct {
	Name  string                 `json:"name"`
	Exams map[string]interface{} `json:"exams"`
}

func GetExamsHash() {
	data, err := ioutil.ReadFile("./ydb/sources/parsed/department_exams_timetable.json")
	newTerms := make(map[string]map[string]interface{})
	examsHash := make(map[string]string)
	if err != nil {
		panic(err)
	}
	var parsedData map[string]map[string]DepartmentExams
	err = json.Unmarshal([]byte(data), &parsedData)
	if err != nil {
		log.Fatal(err)
	}
	for semester, nestedMap := range parsedData {
		newTerms[semester] = make(map[string]interface{})
		newTerms[semester]["exams"] = make(map[string]interface{})
		for departmentLink, nestedValue := range nestedMap {
			if departmentLink != "2603786" {
				continue
			}
			newTerms[semester][departmentLink] = nestedValue.Name
			jsonValue, err := json.MarshalIndent(nestedValue, "", "  ")
			if err != nil {
				panic(err)
			}
			fmt.Printf("Department link: %s\n", departmentLink)
			fmt.Printf("JSON value:\n%s\n", string(jsonValue))
			fmt.Printf("Hash: %x\n", sha256.Sum256(jsonValue))

			for lessonName, lesson := range nestedValue.Exams {
				lessonStruct := parser.LessonData{}
				timeStruct := Time{}
				lessonBytes, err := json.Marshal(lesson)
				if err != nil {
					panic(err)
				}
				err = json.Unmarshal(lessonBytes, &lessonStruct)
				if err != nil {
					panic(err)
				}
				timeBytes, err := json.Marshal(lesson)
				if err != nil {
					panic(err)
				}
				err = json.Unmarshal(timeBytes, &timeStruct)
				if err != nil {
					panic(err)
				}
				lessonDate := lessonStruct.Dates
				lessonDatesFrom := ""
				lessonDatesTo := ""
				weekDay := parser.ConvertDay(strings.Split(lessonDate, ",")[0])
				lessonStruct.Week = weekDay
				//fmt.Println(weekDay)
				lessonDate = strings.Split(lessonDate, ",")[1]
				lessonDate = strings.Split(lessonDate, ".")[0]
				lessonDate = strings.TrimSpace(lessonDate)
				re := regexp.MustCompile(`(\d{2}\s+)([а-яА-Я]+)`)
				lessonDateParts := re.FindStringSubmatch(lessonDate)
				// convert month to number
				month := parser.ConvertMonth(lessonDateParts[2])
				lessonDate = strings.TrimSpace(lessonDateParts[1]) + "." + month
				//fmt.Printf("Lesson dates: %s\n", lessonDate)
				//fmt.Printf("Lesson dates from: %s\n", lessonDatesFrom)
				//fmt.Printf("Lesson dates to: %s\n", lessonDatesTo)
				lessonStruct.DateFrom = lessonDatesFrom
				lessonStruct.DateTo = lessonDatesTo
				lessonStruct.Dates = lessonDate

				lessonTime := timeStruct.Time
				lessonTime = strings.TrimSpace(lessonTime)
				lessonTime = strings.ReplaceAll(lessonTime, " ", "")
				lessonTime = strings.ReplaceAll(lessonTime, " ", "")
				lessonTime = strings.ReplaceAll(lessonTime, "—", "-")
				lessonTimeFrom := strings.Split(lessonTime, "-")[0]
				lessonTimeTo := strings.Split(lessonTime, "-")[1]

				lessonStruct.TimeFrom = lessonTimeFrom
				lessonStruct.TimeTo = lessonTimeTo
				jsonLessonData, err := json.MarshalIndent(lessonStruct, "", "  ")
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(string(jsonLessonData))
				// hash jsonData
				hash := sha256.Sum256([]byte(jsonLessonData))
				hashString := hex.EncodeToString(hash[:])
				// add to examsHash
				examsHash[lessonName] = hashString
				// print lesson as json
				//jsonValue, err := json.MarshalIndent(lessonStruct, "", "  ")
				//if err != nil {
				//	panic(err)
				//}
				//fmt.Printf("JSON value:\n%s\n", string(jsonValue))
				//fmt.Printf("Hash: %x\n", sha256.Sum256(jsonValue))
				// create the original json from lesson struct

			}
			//jsonData, err := json.MarshalIndent(newTerms[semester], "", " ")
			//if err != nil {
			//	fmt.Println(err)
			//}
			// hash jsonData
			//hash := sha256.Sum256([]byte(jsonData))
			//hashString := hex.EncodeToString(hash[:])
			//fmt.Println(hashString)
			//println(string(jsonData))
		}
	}
	jsonData, err := json.MarshalIndent(newTerms, "", " ")
	if err != nil {
		fmt.Println(err)
	}
	// hash jsonData
	hash := sha256.Sum256([]byte(jsonData))
	hashString := hex.EncodeToString(hash[:])
	fmt.Println(hashString)
	//println(string(jsonData))
	fmt.Println(examsHash)
	//fmt.Println(daysHash)
}
