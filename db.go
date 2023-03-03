package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/sugar"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	yc "github.com/ydb-platform/ydb-go-yc"
	"log"
	"os"
	"time"
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

func connect() {
	var cfg Config
	cfg.Endpoint, _ = os.LookupEnv("ENDPOINT")
	cfg.Database, _ = os.LookupEnv("DATABASE")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	db, err := ydb.Open(ctx,
		sugar.DSN(cfg.Endpoint, cfg.Database, true),
		yc.WithInternalCA(),
		yc.WithServiceAccountKeyFileCredentials("./ydb/authorized_key.json"),
	)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	Links, err = getDepartmentLinks(ctx, db)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(len(Links))

	for _, link := range Links {
		fmt.Println(link)
	}

	defer func() {
		_ = db.Close(ctx)
	}()
}
