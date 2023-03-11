package hash

import (
	"fmt"
	"strconv"
	"strings"
)

// oldMap, newMap map[string]string) (newElems, changedElems, missingElems map[string]string
func CompareMaps() {
	oldMap := map[string]string{
		"206771e44e45ee9f4e1943926da310be3da0a5b94b158d743c5bd92289777498": "1",
		"ec052e0cbe9288ac59e1d7224c812908c412280f40f3740c93f345b3acbcefa3": "2",
		"33ebd76541c020c4e710cba48d7371629e460f393f1f78bcd4e477e525d8cbe3": "3",
		"12345": "99",
	}

	newMap := map[string]string{
		"206771e44e45ee9f4e1943926da310be3da0a5b94b158d743c5bd92289777498": "1",
		"33ebd76541c020c4e710cba48d7371629e460f393f1f78bcd4e477e525d8cbe3": "4",
		"d2c195822968b907bea8f090d0daaf53f86e2fca1a1679bfc4afbdba4d4a09e8": "5",
		"123456": "99",
	}

	newElems := make(map[string]string)
	changedKeys := make(map[string]string)
	changedValues := make(map[string]string)
	missingElems := make(map[string]string)

	for key, value := range newMap {
		flag := false
		for oldKey, oldValue := range oldMap {
			if key == oldKey && value != oldValue {
				changedValues[key] = value
				flag = true
				break
			}
			if key != oldKey && value == oldValue {
				changedKeys[key] = value
				flag = true
				break
			}
			if key == oldKey && value == oldValue {
				flag = true
				break
			}
		}
		if !flag {
			newElems[key] = value
		}
	}
	for key, value := range oldMap {
		flag := false
		for newKey, newValue := range newMap {
			if key == newKey || value == newValue {
				flag = true
				break
			}
		}
		if !flag {
			missingElems[key] = value
		}
	}

	fmt.Println("New Elements:", newElems)
	fmt.Println("Changed Elements:", changedKeys)
	fmt.Println("Elements with Changed IDs:", changedValues)
	fmt.Println("Missing Elements:", missingElems)
}

func compareDates(date1 string, date2 string) int {
	day1, _ := strconv.Atoi(strings.Split(date1, ".")[0])
	month1, _ := strconv.Atoi(strings.Split(date1, ".")[1])
	day2, _ := strconv.Atoi(strings.Split(date2, ".")[0])
	month2, _ := strconv.Atoi(strings.Split(date2, ".")[1])

	if month1 < month2 {
		return -1
	} else if month1 > month2 {
		return 1
	} else {
		if day1 < day2 {
			return -1
		} else if day1 > day2 {
			return 1
		} else {
			return 0
		}
	}
}
