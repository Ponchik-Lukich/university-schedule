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

type Department struct {
	Id   string
	Name string
}

type Config struct {
	Endpoint string
	Database string
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

//func createTableExample(ctx context.Context, db ydb.Connection, cfg Config) {
//	err := db.Table().Do(ctx,
//		func(ctx context.Context, s table.Session) (err error) {
//			return s.CreateTable(ctx, cfg.Database+"/series",
//				options.WithColumn("series_id", types.TypeUint64), // not null column
//				options.WithColumn("title", types.Optional(types.TypeUTF8)),
//				options.WithColumn("release_date", types.Optional(types.TypeDate)),
//				options.WithPrimaryKeyColumn("series_id"),
//			)
//		},
//	)
//	if err != nil {
//		panic(err)
//	}
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
	err = db.Table().Do(ctx, func(ctx context.Context, s table.Session) error {
		query := `SELECT * FROM departments`
		_, res, err := s.Execute(ctx, table.DefaultTxControl(), query, table.NewQueryParameters())
		if err != nil {
			return err
		}
		fmt.Println(res.ResultSetCount())
		for res.NextResultSet(ctx) {
			for res.NextRow() {
				department := &Department{}
				err := res.ScanWithDefaults(
					&department.Id,
					&department.Name,
				)
				if err != nil {
					return err
				}
				fmt.Println(*department)
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer func() {
		_ = db.Close(ctx)
	}()
}
