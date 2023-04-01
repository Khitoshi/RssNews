package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
)

type USER struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	Id         int64     `bun:"id,pk,autoincrement"`
	Name       string    `bun:"name,notnull"`
	Email      string    `bun:"email,notnull,unique"`
	Password   string    `bun:"password,notnull"`
	Created_at time.Time `bun:"created_at,notnull"`
	Updated_at time.Time `bun:"updated_at,notnull"`
}

type ITEMS struct {
	bun.BaseModel `bun:"table:items,alias:i"`

	Id           int64     `bun:"id,pk,autoincrement"`
	Url          string    `bun:"url,notnull,unique"`
	Title        string    `bun:"title,notnull"`
	Description  string    `bun:"description,notnull"`
	Author       string    `bun:"author"`
	Published_at time.Time `bun:"published_at"`
	Created_at   time.Time `bun:"created_at,notnull"`
	Updated_at   time.Time `bun:"updated_at,notnull"`
}

type USER_ITEMS struct {
	bun.BaseModel `bun:"table:user_items,alias:ui"`

	User_id    int64     `bun:"user_id,FOREIGN KEY"`
	Article_id int64     `bun:"article_id,FOREIGN KEY"`
	Created_at time.Time `bun:"created_at,not null"`
}

func main() {
	startTime := time.Now()
	fmt.Printf("start time: %v \n", startTime)
	defer func() {
		fmt.Printf("\n processing time: %v", time.Since(startTime).Milliseconds())
	}()

	//dbを開く
	sqldb, err := sql.Open("postgres", "user=postgres dbname=rss_reader_web password=985632 sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer sqldb.Close()

	db := bun.NewDB(sqldb, pgdialect.New())
	defer db.Close()
	//クエリのパラメーター出力
	db.AddQueryHook(bundebug.NewQueryHook(
		//bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))

	ctx := context.Background()
	_, err = db.NewCreateTable().Model((*USER)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		log.Fatal(err)
	}

	ctx = context.Background()
	_, err = db.NewCreateTable().Model((*ITEMS)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		log.Fatal(err)
	}

}
