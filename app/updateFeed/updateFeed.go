package updateFeed

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	tables "rss_reader/tables"

	"time"

	_ "github.com/lib/pq"
	"github.com/mmcdole/gofeed"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
)

// itemの定期更新
func UpdateItemsFromRSSFeed() error {
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

	//一時間の周期処理
	//ticker := time.NewTicker(1 * time.Hour)
	ticker := time.NewTicker(10 * time.Second)
	for range ticker.C {

		fmt.Println("定期処理実行")

		//RSSから情報取得
		rssfeed, err := loadRssFeed(db)
		if err != nil {
			//log.Fatal(err)
			return err
		}

		for _, feed := range rssfeed {
			err = update(db, &feed)
			if err != nil {
				//log.Fatal(err)
				return err
			}
		}

	}

	return nil
}

// 記事更新処理
func update(db *bun.DB, feeds *gofeed.Feed) error {
	for _, item := range feeds.Items {
		fmt.Println("周期処理スタート:", time.Now())

		f := tables.ITEMS{}
		//すでに同じリンクが存在しないかチェック
		err := db.NewSelect().Model(&f).Where("url=?", item.Link).Scan(context.Background())
		if err != nil {
			return err
		}
		//URLが空の場合,重複なし
		if f.Url != "" {
			fmt.Println("break")
			break
		}

		//INSERT処理
		f = tables.ITEMS{
			//Id:           nil,
			Url:          item.Link,
			Title:        item.Title,
			Description:  item.Description,
			Author:       item.Author.Name,
			Published_at: *item.PublishedParsed,
			Created_at:   time.Now(),
			Updated_at:   *item.UpdatedParsed,
		}
		_, err = db.NewInsert().Model(&f).Exec(context.Background())
		if err != nil {
			return err
		}
		fmt.Printf("insert %v\n", f.Title)
	}

	return nil
}

func loadRssFeed(db *bun.DB) ([]gofeed.Feed, error) {

	rssURLs := []tables.RSS_URLS{}
	//すでに同じリンクが存在しないかチェック
	err := db.NewSelect().Model(&rssURLs).Where("rss_id = 0").Scan(context.Background())
	if err != nil {
		return nil, err
	}
	rssFeed := []gofeed.Feed{}
	//rssURLs := []string{}
	//rssURLs = append(rssURLs, "https://qiita.com/IXKGAGB/feed")

	for _, url := range rssURLs {
		rf, err := gofeed.NewParser().ParseURL(url.Rss_URL)
		if err != nil {
			return nil, err
		}
		rssFeed = append(rssFeed, *rf)
	}

	return rssFeed, nil

}
