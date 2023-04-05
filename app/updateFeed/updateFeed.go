package updateFeed

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	items "rss_reader/tables"
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

	//TODO:定期処理を10秒に設定してあるので1時間に変更
	ticker := time.NewTicker(10 * time.Second)
	for range ticker.C {

		fmt.Println("定期処理実行")

		//RSSから情報取得
		rssfeed, err := gofeed.NewParser().ParseURL("https://qiita.com/IXKGAGB/feed")
		if err != nil {
			//log.Fatal(err)
			return err
		}

		//更新
		for _, item := range rssfeed.Items {
			f := items.ITEMS{}
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
			f = items.ITEMS{
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

	}

	return nil
}
