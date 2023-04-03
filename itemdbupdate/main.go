package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"

	_ "github.com/lib/pq"
)

func updateFeed() error {
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

	rssfeed, err := gofeed.NewParser().ParseURL("https://qiita.com/IXKGAGB/feed")
	if err != nil {
		//log.Fatal(err)
		return err
	}

	stmt, err := db.Prepare("INSERT INTO items(url,title,description,author,published_at,created_at, updated_at) VALUES($1,$2,$3,$4,$5,$6,$7);")
	if err != nil {
		return err
	}
	for _, item := range rssfeed.Items {

		fmt.Println("update:", item.Link)
		res, err := stmt.Exec(
			item.Link,
			item.Title,
			item.Description,
			item.Author,
			item.PublishedParsed,
			time.Now(),
			item.UpdatedParsed,
		)
		if err != nil {
			return nil
		}

	}

	return nil
}

func main() {
	startTime := time.Now()
	fmt.Printf("start time: %v \n", startTime)
	defer func() {
		fmt.Printf("\n processing time: %v", time.Since(startTime).Milliseconds())
	}()

	err := updateFeed()
	if err != nil {
		log.Fatal(err)
	}
}
