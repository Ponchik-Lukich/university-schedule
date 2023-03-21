package main

import (
	"sync"
	"university-timetable/parser"
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
	parser.CreateInserts()
	//var cfg Config
	//cfg.Endpoint, _ = os.LookupEnv("ENDPOINT")
	//cfg.Database, _ = os.LookupEnv("DATABASE")
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	//defer cancel()
	//db, err := ydb.Open(ctx,
	//	sugar.DSN(cfg.Endpoint, cfg.Database, true),
	//	yc.WithInternalCA(),
	//	yc.WithServiceAccountKeyFileCredentials("./ydb/authorized_key.json"),
	//)
	//if err != nil {
	//	fmt.Println(err)
	//	panic(err)
	//}
	//Links, err = getDepartmentLinks(ctx, db, 16)
	//if err != nil {
	//	fmt.Println(err)
	//	panic(err)
	//}
	//for _, link := range Links {
	//	fmt.Println(link)
	//}

	//LessonsData, err := getLessonsData(ctx, db, 2603786)
	//if err != nil {
	//	fmt.Println(err)
	//	panic(err)
	//}
	//for _, lesson := range LessonsData {
	//	fmt.Println(lesson)
	//}

	//RoomsData, err := getRooms(ctx, db)
	//for _, room := range RoomsData {
	//	fmt.Println(room)
	//}
	//
	//defer func() {
	//	_ = db.Close(ctx)
	//}()

	//hash.GetDepartmentsHash()
	//hash.GetExamsHash()
	//hash.CompareMaps()
	//parser.ParseByXpath("https://home.potatohd.ru/departments/2603786")
	//parser.ParseByXpathExam("https://home.potatohd.ru/departments/2603786/exams")
	//parser.ParseRoomByXpath("https://home.mephi.ru/rooms/4711947")

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
