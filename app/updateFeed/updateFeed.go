package updateFeed

import (
	"context"
	"database/sql"
	"fmt"
	"rss_reader/database"
	"rss_reader/tables"

	"time"

	_ "github.com/lib/pq"
	"github.com/mmcdole/gofeed"
	"github.com/uptrace/bun"
)

// 定期更新用
func RegularUpdatingOfArticles() error {
	//WARNING:デバッグ用に10秒に1回更新をしている可能性があるので要チェック!
	ticker := time.NewTicker(1 * time.Hour)
	//ticker := time.NewTicker(10 * time.Second)
	for range ticker.C {
		fmt.Println("定期処理実行開始:", time.Now())
		//更新処理
		err := UpdateItemsFromRSSFeed()
		if err != nil {
			return err
		}
	}
	return nil
}

// 記事の更新処理
func UpdateItemsFromRSSFeed() error {

	//一時間の周期処理
	//rssのURL群を取得
	rssURLs, err := getRssURL()
	if err != nil {
		return err
	}

	for _, rssURL := range rssURLs {
		err = RegisterRSSFeeds(rssURL)
		if err != nil {
			return err
		}
	}

	return nil
}

func RegisterRSSFeeds(rssURL tables.RSS_URLS) error {
	//URLから記事群を取得
	feeds, err := getFeed(rssURL.Rss_URL)
	if err != nil {
		return err
	}

	//記事群をテーブルにinsert
	err = insertFeeds(feeds, rssURL.Rss_id)
	if err != nil {
		return err
	}
	return nil
}

// テーブルからrssURLを取得
func getRssURL() ([]tables.RSS_URLS, error) {
	rssURLs := []tables.RSS_URLS{}
	database.WithDBConnection(func(db *bun.DB) error {
		err := db.NewSelect().Model(&rssURLs).Scan(context.Background())
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		return nil
	})
	return rssURLs, nil
}

// URLから記事を取得
func getFeed(rssurl string) (*gofeed.Feed, error) {
	//rssurlsから記事群を取得
	f, err := gofeed.NewParser().ParseURL(rssurl)
	//f, err := gofeed.NewParser().ParseURL("https://qiita.com/popular-items/feed.atom")
	if err != nil {
		return nil, err
	}

	return f, nil
}

// 記事がすでに登録されているかのチェック trueの場合存在する
func isFeedExist(url string) (bool, error) {
	feed := tables.ITEMS{}
	err := database.WithDBConnection(func(db *bun.DB) error {
		err := db.NewSelect().Model(&feed).Where("url = ?", url).Scan(context.Background())
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		return nil
	})
	if err != nil {
		return false, err
	}

	//存在する場合
	if feed.Url != "" {
		return true, nil
	}
	//存在しなかった場合
	return false, nil
}

// 記事群をテーブルにINSERT
func insertFeeds(feeds *gofeed.Feed, rssID int64) error {
	//登録処理
	for _, feed := range feeds.Items {

		//テーブルに同じ記事が存在しないかチェック
		isFeed, err := isFeedExist(feed.Link)
		if err != nil {
			return err
		}
		//存在する場合スキップする
		if isFeed {
			continue
		}

		//テーブルのカラムの構造体に置き換え
		item := tables.ITEMS{
			Rss_id:       rssID,
			Title:        feed.Title,
			Url:          feed.Link,
			Description:  feed.Description,
			Author:       feed.Author.Name,
			Published_at: time.Now(),
			Created_at:   time.Now(),
			Updated_at:   time.Now(),
		}

		//テーブルにINSERT
		err = database.WithDBConnection(func(db *bun.DB) error {
			err := db.NewInsert().Model(&item).Scan(context.Background())
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
		fmt.Println("INSERT 記事:", feed.Title)
	}

	return nil
}
