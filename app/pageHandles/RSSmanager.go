package pageHandles

import (
	"context"
	"database/sql"
	"net/http"
	"rss_reader/database"
	"rss_reader/tables"
	"rss_reader/updateFeed"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
)

func HandleRSSList_Get(c echo.Context) error {

	//クッキーからユーザーID取得
	cookie, err := c.Cookie("userId")
	if err != nil {
		return c.Redirect(http.StatusFound, "login")
	}
	//クッキーから入手したuseridをintに変換
	userID, err := strconv.ParseInt(cookie.Value, 10, 64)
	if err != nil {
		//return err
		return c.Redirect(http.StatusFound, "login")
	}

	//useritemsテーブル ユーザーIDからrssIDを取得
	rssID := []tables.USER_ITEMS{}
	err = database.WithDBConnection(func(db *bun.DB) error {
		err := db.NewSelect().Model(&rssID).Column("rss_id").Where("user_id=?", userID).Scan(context.Background())
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	//rssIDからrssURLを取得
	rssURLs := make([]tables.RSS_URLS, len(rssID))
	for i, id := range rssID {
		err = database.WithDBConnection(func(db *bun.DB) error {
			err := db.NewSelect().Model(&rssURLs[i]).Where("rss_id=?", id.Rss_id).Scan(context.Background())
			if err != nil && err != sql.ErrNoRows {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	return c.Render(http.StatusOK, "feedlist", rssURLs)
}

func HandleRSSList_Post(c echo.Context) error {
	//クッキーからユーザーID取得
	cookie, err := c.Cookie("userId")
	if err != nil {
		return c.Redirect(http.StatusFound, "login")
	}
	//クッキーから入手したuseridをintに変換
	userID, err := strconv.ParseInt(cookie.Value, 10, 64)
	if err != nil {
		//return err
		return c.Redirect(http.StatusFound, "login")
	}

	//input値を取得
	for "" != c.FormValue("rssURL") {
		//TODO:ずっとzennのURLが入っている状態なのでチェックする
		rssURL := c.FormValue("rssURL")

		//すでに同じURLが登録されていないかチェック
		Rssid, err := hasRSSURLAlreadyRegistered(rssURL)
		if err != nil {
			return err
		}

		//存在しなかった時,rssURLsに登録する
		if Rssid == -1 {
			//rssurlsに登録
			registerRSS(rssURL)

			//rssURL rssfeedsを入手してテーブルに入れる

			//すでに同じURLが登録されていないかチェック
			Rssid, err = hasRSSURLAlreadyRegistered(rssURL)
			if err != nil || Rssid == -1 {
				return err
			}

			rss := tables.RSS_URLS{
				Rss_id:  Rssid,
				Rss_URL: rssURL,
			}

			//追加したrssの記事をテーブルに追加
			updateFeed.RegisterRSSFeeds(rss)
		}
		//user_itemに登録
		RegisterSubscriptionUserItem(Rssid, userID)
	}

	return c.Redirect(http.StatusFound, "/")
}
