package loadFeed

import (
	"context"
	"rss_reader/database"
	"rss_reader/tables"
	table_items "rss_reader/tables"

	//"rss_reader/updateFeed"
	"time"

	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
)

type Feed struct {
	Title           string //タイトル
	Link            string //記事リンク
	Description     string
	PublishedParsed *time.Time
	UpdatedParsed   *time.Time
}

/*
// テーブルからitemを取得
func GetFeeds(userId int64) ([]News, error) {

	//SELECT * FROM items  LEFT JOIN  (SELECT * FROM user_items WHERE user_id = 8) AS rssid ON items.rss_id = rssid.rss_id;
	//このsqlをbunに変換する

	user_items := []tables.USER_ITEMS{}
	//err = db.NewSelect().Model(user_items).Where("user_id=?", 8).Scan(context.Background())
	//err = db.NewSelect().Model(user_items).Column("rss_id").Where("user_id = ?", 8).Scan(context.Background())

	//rssid取得
	err := database.WithDBConnection(func(db *bun.DB) error {
		err := db.NewSelect().Model(&user_items).Column("rss_id").Where("user_id = ?", userId).Scan(context.Background())
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	news := []News{}

	for i, item := range user_items {
		feed := []Feed{}
		items := []table_items.ITEMS{}

		//
		err := database.WithDBConnection(func(db *bun.DB) error {
			err = db.NewSelect().Model(&items).Where("rss_id = ?", item.Rss_id).Scan(context.Background())
			if err != nil {
				return err
			}
			return nil
		})
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
		n := News{
			siteName: fmt.Sprintf("%v", i),
			feeds:    feed,
		}
		news = append(news, n)
	}

	return news, nil
}
*/

// テーブルからitemを取得
func GetFeeds(userId int64) ([]Feed, error) {

	//SELECT * FROM items  LEFT JOIN  (SELECT * FROM user_items WHERE user_id = 8) AS rssid ON items.rss_id = rssid.rss_id;
	//このsqlをbunに変換する
	user_items := []tables.USER_ITEMS{}

	//rssid取得
	err := database.WithDBConnection(func(db *bun.DB) error {
		err := db.NewSelect().Model(&user_items).Column("rss_id").Where("user_id = ?", userId).Scan(context.Background())
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	feed := []Feed{}
	for _, item := range user_items {
		items := []table_items.ITEMS{}
		//
		err := database.WithDBConnection(func(db *bun.DB) error {
			err = db.NewSelect().Model(&items).Where("rss_id = ?", item.Rss_id).Limit(0).Limit(9).Scan(context.Background())
			if err != nil {
				return err
			}
			return nil
		})
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
