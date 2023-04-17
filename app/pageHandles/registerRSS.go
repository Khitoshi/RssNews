package pageHandles

import (
	"context"
	"database/sql"
	"net/http"
	"rss_reader/database"
	"rss_reader/tables"
	"rss_reader/updateFeed"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
)

func HandleRegisterRSS_Get(c echo.Context) error {
	return c.Render(http.StatusOK, "registerRSS", nil)

}

func HandleRegisterRSS_Post(c echo.Context) error {
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

	return c.Redirect(http.StatusFound, "/")
}

// すでに同じURLが存在しないかテーブルをチェック
// 存在しなかった場合-1が返る
func hasRSSURLAlreadyRegistered(url string) (int64, error) {
	r := tables.RSS_URLS{}
	err := database.WithDBConnection(func(db *bun.DB) error {
		err := db.NewSelect().Model(&r).Where("rss_url=?", url).Scan(context.Background())
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		return nil
	})
	if err != nil {
		return -1, err
	}

	//検索したテーブルにデータが存在する場合 rssIDを返す
	if r.Rss_URL != "" {
		return r.Rss_id, nil
	}

	//存在しなかったので-1
	return -1, nil
}

// rssURLの登録
func registerRSS(rssurl string) error {

	r := tables.RSS_URLS{

		Rss_URL:    rssurl,
		Created_at: time.Now(),
	}

	err := database.WithDBConnection(func(db *bun.DB) error {
		_, err := db.NewInsert().Model(&r).Exec(context.Background())
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// userItemに既に登録されているかのチェック//return falseの時、登録されていない
func HasRegisterSubscriptionUserItem(rssID, userID int64) (bool, error) {
	//useritemsを追加
	u := tables.USER_ITEMS{}
	err := database.WithDBConnection(func(db *bun.DB) error {
		err := db.NewSelect().Model(&u).Where("rss_id=? AND user_id=?", rssID, userID).Scan(context.Background())
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		return nil
	})
	if err != nil {
		return false, err
	}

	if u.Rss_id == 0 {
		return false, nil
	}

	return true, nil
}

// userItemに登録
func RegisterSubscriptionUserItem(rssID, userID int64) error {
	//useritemsを追加
	u := tables.USER_ITEMS{
		Rss_id:     rssID,
		User_id:    userID,
		Created_at: time.Now(),
	}
	err := database.WithDBConnection(func(db *bun.DB) error {
		_, err := db.NewInsert().Model(&u).Exec(context.Background())
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil

}
