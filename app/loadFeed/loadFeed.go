package loadFeed

import (
	"context"
	"database/sql"
	"log"
	table_items "rss_reader/tables"

	//"rss_reader/updateFeed"
	"time"

	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
)

type Feed struct {
	Title           string //タイトル
	Link            string //記事リンク
	Description     string
	PublishedParsed *time.Time
	UpdatedParsed   *time.Time
}

// テーブルからitemを取得
func GetFeeds() ([]Feed, error) {
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

	items := []table_items.ITEMS{}
	err = db.NewSelect().Model(&items).Scan(context.Background())
	if err != nil {
		return nil, err
	}

	f := make([]Feed, len(items))
	for i, item := range items {
		f[i] = Feed{
			Title:           item.Title,
			Link:            item.Url,
			Description:     item.Description,
			PublishedParsed: &item.Published_at,
			UpdatedParsed:   &item.Updated_at,
		}
	}
	return f, nil
}
