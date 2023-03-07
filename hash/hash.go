package hash

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type Department struct {
	Name string                 `json:"name"`
	Days map[string]interface{} `json:"days"`
}

type Lesson struct {
	Name           string            `json:"name"`
	Time           string            `json:"time"`
	Dates          string            `json:"dates"`
	Room           string            `json:"room"`
	RoomID         string            `json:"room_id"`
	Week           string            `json:"week"`
	Type           string            `json:"type"`
	Groups         map[string]string `json:"groups"`
	Tutors         map[string]string `json:"tutors"`
	AdditionalInfo []string          `json:"additional_info"`
}

func GetHash() {
	data, err := ioutil.ReadFile("./ydb/sources/parsed/department_timetable.json")
	if err != nil {
		panic(err)
	}

	var parsedData map[string]map[string]Department
	//var parsedLessons map[string]map[string]map[string]Lesson

	err = json.Unmarshal([]byte(data), &parsedData)
	if err != nil {
		log.Fatal(err)
	}
	for semester, nestedMap := range parsedData {
		fmt.Printf("Semester: %s\n", semester)
		for departmentLink, nestedValue := range nestedMap {
			jsonValue, err := json.MarshalIndent(nestedValue, "", "  ")
			if err != nil {
				panic(err)
			}
			fmt.Printf("Department link: %s\n", departmentLink)
			//fmt.Printf("JSON value:\n%s\n", string(jsonValue))
			fmt.Printf("Hash: %x\n", sha256.Sum256(jsonValue))
			// print lessons
			for day, lessons := range nestedValue.Days {
				fmt.Printf("Day: %s\n", day)
				for lessonName, lesson := range lessons.(map[string]interface{}) {
					fmt.Printf("Lesson name: %s\n", lessonName)
					jsonValue, err := json.MarshalIndent(lesson, "", "  ")
					if err != nil {
						panic(err)
					}
					fmt.Println(string(jsonValue))
					fmt.Printf("Hash: %x\n", sha256.Sum256(jsonValue))
				}
			}
		}
	}

}
