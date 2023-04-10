package database

import (
	"context"
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
)

func WithDBConnection(f func(db *bun.DB) error) error {
	//dbを開く
	sqldb, err := sql.Open("postgres", "user=postgres dbname=rss_reader_web password=985632 sslmode=disable")
	if err != nil {
		//log.Fatal(err)
		return err
	}
	defer sqldb.Close()

	db := bun.NewDB(sqldb, pgdialect.New())
	defer db.Close()

	//クエリのパラメーター出力
	db.AddQueryHook(bundebug.NewQueryHook(
		//bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))

	err = f(db)
	if err != nil {
		//log.Fatal(err)
		return err
	}

	return nil
}

// テーブル作成
// model:テーブル構造体
func CreateTable(model any) error {
	err := WithDBConnection(func(db *bun.DB) error {
		ctx := context.Background()
		_, err := db.NewCreateTable().Model(model).IfNotExists().Exec(ctx)

		if err != nil {
			//log.Fatal(err)
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
