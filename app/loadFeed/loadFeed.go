package loadFeed

import (
	"context"
	"database/sql"
	"log"
	"rss_reader/tables"
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

	//SELECT * FROM items  LEFT JOIN  (SELECT * FROM user_items WHERE user_id = 8) AS rssid ON items.rss_id = rssid.rss_id;
	//このsqlをbunに変換する
	//TODO: coocieからuseridを入手に変更
	user_items := []tables.USER_ITEMS{}
	//err = db.NewSelect().Model(user_items).Where("user_id=?", 8).Scan(context.Background())
	//err = db.NewSelect().Model(user_items).Column("rss_id").Where("user_id = ?", 8).Scan(context.Background())

	//rssid取得
	err = db.NewSelect().Model(&user_items).Column("rss_id").Where("user_id = ?", 8).Scan(context.Background())
	if err != nil {
		return nil, err
	}

	feed := []Feed{}
	for _, item := range user_items {
		items := []table_items.ITEMS{}
		err = db.NewSelect().Model(&items).Where("rss_id = ?", item.Rss_id).Scan(context.Background())
		if err != nil {
			return nil, err
		}

		for _, item := range items {
			f := Feed{
				Title:           item.Title,
				Link:            item.Url,
				Description:     item.Description,
				PublishedParsed: &item.Published_at,
				UpdatedParsed:   &item.Updated_at,
			}
			feed = append(feed, f)
		}
	}

	return feed, nil
}
