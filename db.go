package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"log"
)

type DepartmentLink struct {
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

func getDepartmentLinks(ctx context.Context, db ydb.Connection) ([]string, error) {
	var departmentLinks []string
	err := db.Table().Do(ctx, func(ctx context.Context, s table.Session) error {
		query := `SELECT id FROM department_links
				  WHERE semester=16;`
		_, res, err := s.Execute(ctx, table.DefaultTxControl(), query, table.NewQueryParameters())
		if err != nil {
			return err
		}
		for res.NextResultSet(ctx) {
			for res.NextRow() {
				department := &DepartmentLink{}
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

//func getLessonsData(ctx context.Context, db ydb.Connection, departmentLink int) ([]string, error) {
//
//}
