package updateFeed

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"rss_reader/database"
	tables "rss_reader/tables"

	"time"

	_ "github.com/lib/pq"
	"github.com/mmcdole/gofeed"
	"github.com/uptrace/bun"
)

// itemの定期更新
func UpdateItemsFromRSSFeed() error {
	//一時間の周期処理
	//ticker := time.NewTicker(10 * time.Second)
	ticker := time.NewTicker(1 * time.Hour)
	for range ticker.C {

		fmt.Println("定期処理実行")

		//RSSから情報取得
		rssfeed, err := loadRssFeed()
		if err != nil {
			//log.Fatal(err)
			return err
		}

		for _, feed := range rssfeed {
			err = update(&feed)
			if err != nil {
				//log.Fatal(err)
				log.Fatal(err)
				return err
			}
		}

	}

	return nil
}

// 記事更新処理
func update(feeds *gofeed.Feed) error {
	for _, item := range feeds.Items {
		fmt.Println("周期処理スタート:", time.Now())

		f := tables.ITEMS{}
		//すでに同じリンクが存在しないかチェック
		err := database.WithDBConnection(func(db *bun.DB) error {
			err := db.NewSelect().Model(&f).Where("url=?", item.Link).Scan(context.Background())
			if err != nil && err != sql.ErrNoRows {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}

		//URLが空の場合,重複なし
		if f.Url != "" {
			fmt.Println("break")
			break
		}

		//TODO:zennの場合 UpdatedParsed が実装されていない
		//サイトによって実装されていない項目もあるので各項目でnullchekして値がない場合の処理を作る
		//INSERT処理
		f = tables.ITEMS{
			//Id:           nil,
			Url:          item.Link,
			Title:        item.Title,
			Description:  item.Description,
			Author:       item.Author.Name,
			Published_at: *item.PublishedParsed,
			Created_at:   time.Now(),
			Updated_at:   time.Now(),
			//Updated_at:   *item.UpdatedParsed,
		}
		//テーブルに記事を登録
		err = database.WithDBConnection(func(db *bun.DB) error {
			_, err = db.NewInsert().Model(&f).Exec(context.Background())
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return nil
		}
		fmt.Printf("insert %v\n", f.Title)
	}

	return nil
}

func loadRssFeed() ([]gofeed.Feed, error) {

	rssURLs := []tables.RSS_URLS{}
	//すでに同じリンクが存在しないかチェック
	err := database.WithDBConnection(func(db *bun.DB) error {
		err := db.NewSelect().Model(&rssURLs).Scan(context.Background())
		if err != nil {
			return err
		}
		return nil
	})
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
