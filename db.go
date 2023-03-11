package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"log"
	"strconv"
	"university-timetable/parser"
)

type ID struct {
	Id uint64
}

type Config struct {
	Endpoint string
	Database string
}

var Links []string

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func getDepartmentLinks(ctx context.Context, db ydb.Connection, semester int) ([]string, error) {
	var departmentLinks []string
	err := db.Table().Do(ctx, func(ctx context.Context, s table.Session) error {
		query := fmt.Sprintf(`SELECT id FROM department_links
				  WHERE semester=%d;`, semester)
		_, res, err := s.Execute(ctx, table.DefaultTxControl(), query, table.NewQueryParameters())
		if err != nil {
			return err
		}
		for res.NextResultSet(ctx) {
			for res.NextRow() {
				department := &ID{}
				err := res.ScanWithDefaults(
					&department.Id,
				)
				if err != nil {
					return err
				}
				departmentLinks = append(departmentLinks, fmt.Sprintf("%d", department.Id))
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return departmentLinks, nil
}

func getLessonsData(ctx context.Context, db ydb.Connection, departmentLink int) (map[string]parser.LessonData, error) {
	var lessonsData map[string]parser.LessonData
	err := db.Table().Do(ctx, func(ctx context.Context, s table.Session) error {
		query := fmt.Sprintf(`SELECT cp.id AS id, cp.time_from AS time_from, cp.time_to AS time_to, cp.type AS type, 
cp.week AS week, cp.subject AS name, cp.week_day AS week_day, cp.room_id AS room_id, cp.date_from as date_from,
cp.date_to AS date_to, cp.date_from AS dates FROM calendar_plan cp
INNER JOIN calendar_plan_department_links cpdl on cp.id = cpdl.calendar_plan_id
INNER JOIN department_links dl on cpdl.department_link_id = dl.id
WHERE dl.id = %d;`, departmentLink)
		_, res, err := s.Execute(ctx, table.DefaultTxControl(), query, table.NewQueryParameters())
		if err != nil {
			return err
		}
		for res.NextResultSet(ctx) {
			for res.NextRow() {
				lessonName := &ID{}
				lesson := &parser.LessonData{}
				err := res.ScanWithDefaults(
					&lessonName.Id,
					&lesson.TimeFrom,
					&lesson.TimeTo,
					&lesson.Type,
					&lesson.Week,

					&lesson.Name,

					//&lesson.Tutors,
					//&lesson.Groups,
					&lesson.Room,
					&lesson.RoomID,
					&lesson.DateFrom,
					&lesson.DateTo,
					&lesson.Dates,
				)
				if err != nil {
					return err
				}
				lessonsData[strconv.Itoa(int(lessonName.Id))] = *lesson
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return lessonsData, nil
}
