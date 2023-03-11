package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"university-timetable/parser"
)

type Department struct {
	Name string                 `json:"name"`
	Days map[string]interface{} `json:"days"`
}

type Time struct {
	Time string `json:"time"`
}

func GetDepartmentsHash() {
	data, err := ioutil.ReadFile("./ydb/sources/parsed/department_timetable.json")
	newTerms := make(map[string]map[string]interface{})
	daysHash := make(map[string]string)
	lessonsHash := make(map[string]string)
	if err != nil {
		panic(err)
	}
	var parsedData map[string]map[string]Department

	err = json.Unmarshal([]byte(data), &parsedData)
	if err != nil {
		log.Fatal(err)
	}
	//counter := 0
	for semester, nestedMap := range parsedData {
		newTerms[semester] = make(map[string]interface{})
		newTerms[semester]["days"] = make(map[string]interface{})
		for departmentLink, nestedValue := range nestedMap {
			if departmentLink != "2604078" {
				continue
			}
			newTerms[semester][departmentLink] = nestedValue.Name
			//jsonValue, err := json.MarshalIndent(nestedValue, "", "  ")
			//if err != nil {
			//	panic(err)
			//}
			//fmt.Printf("Department link: %s\n", departmentLink)
			//fmt.Printf("JSON value:\n%s\n", string(jsonValue))
			//fmt.Printf("Hash: %x\n", sha256.Sum256(jsonValue))

			// print lessons
			for day, lessons := range nestedValue.Days {
				//fmt.Printf("Day: %s\n", day)
				dayData := make(map[string]parser.LessonData)
				for lessonName, lesson := range lessons.(map[string]interface{}) {
					//fmt.Printf("Lesson name: %s\n", lessonName)
					// parse lesson in struct
					fmt.Println("NAME: ", lessonName)
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
					if lessonDate != "" {
						lessonDates := strings.TrimSpace(lessonDate)
						lessonDates = strings.ReplaceAll(lessonDates, " ", "")
						lessonDates = strings.ReplaceAll(lessonDates, " ", "")
						if strings.ContainsAny(lessonDates, "-") {
							lessonDates = strings.ReplaceAll(lessonDates, "(", "")
							lessonDates = strings.ReplaceAll(lessonDates, ")", "")
							lessonDates = strings.ReplaceAll(lessonDates, ",", "-")
							lessonDatesFrom = strings.Split(lessonDates, "-")[0]
							lessonDatesTo = strings.Split(lessonDates, "-")[1]
							lessonDate = ""
						} else {
							lessonDates = strings.ReplaceAll(lessonDates, "(", "")
							lessonDates = strings.ReplaceAll(lessonDates, ")", "")
							lessonDate = lessonDates
						}
					}
					//fmt.Printf("Lesson dates: %s\n", lessonDate)
					//fmt.Printf("Lesson dates from: %s\n", lessonDatesFrom)
					//fmt.Printf("Lesson dates to: %s\n", lessonDatesTo)
					lessonStruct.DateFrom = lessonDatesFrom
					lessonStruct.DateTo = lessonDatesTo
					lessonStruct.Dates = lessonDate
					//lessonStruct.Week = parser.ConvertDay(day)
					lessonStruct.WeekDay = parser.ConvertDay(day)

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
					// add to lessonsHash
					lessonsHash[hashString] = lessonName
					dayData[lessonName] = lessonStruct

					// print lesson as json
					//jsonValue, err := json.MarshalIndent(lessonStruct, "", "  ")
					//if err != nil {
					//	panic(err)
					//}
					//fmt.Printf("JSON value:\n%s\n", string(jsonValue))
					//fmt.Printf("Hash: %x\n", sha256.Sum256(jsonValue))
					// create the original json from lesson struct

				}
				newTerms[semester]["days"].(map[string]interface{})[day] = dayData
				jsondayData, err := json.MarshalIndent(dayData, "", " ")
				if err != nil {
					fmt.Println(err)
				}
				// hash jsonData
				hash := sha256.Sum256([]byte(jsondayData))
				hashString := hex.EncodeToString(hash[:])
				// add to daysHash
				daysHash[hashString] = day
			}
			jsonData, err := json.MarshalIndent(newTerms[semester], "", " ")
			if err != nil {
				fmt.Println(err)
			}
			// hash jsonData
			hash := sha256.Sum256([]byte(jsonData))
			hashString := hex.EncodeToString(hash[:])
			fmt.Println(hashString)
			//println(string(jsonData))
		}
	}
	fmt.Println(lessonsHash)
	fmt.Println(daysHash)
	//jsonData, err := json.MarshalIndent(newTerms, "", " ")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//// hash jsonData
	//hash := sha256.Sum256([]byte(jsonData))
	//hashString := hex.EncodeToString(hash[:])
	//fmt.Println(hashString)
	//fmt.Println(lessonsHash)
	//fmt.Println(daysHash)
	//println(string(jsonData))

}
