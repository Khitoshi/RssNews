package pageHandles

import (
	"context"
	"database/sql"
	"net/http"
	"rss_reader/database"
	"rss_reader/tables"

	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
)

func HandleRegisterRSS_Get(c echo.Context) error {
	return c.Render(http.StatusOK, "registerRSS", nil)

}

func HandleRegisterRSS_Post(c echo.Context) error {
	rssURL := c.FormValue("rssURL")

	hasRssURL,err := IsRSSURLAlreadyRegistered(rssURL)
	if err != nil {
		return err
	}

	//
	if hasRssURL {//存在しなかった時の処理
		registerRSS(&r)
		
	}
	else{//存在した時の処理
		
	}


	return c.Redirect(http.StatusFound, "/")
}

// すでに同じURLが存在しないかをチェック falseの場合存在していた
func IsRSSURLAlreadyRegistered(url string) (bool, error) {
	r := tables.RSS_URLS{}
	err := database.WithDBConnection(func(db *bun.DB) error {
		err := db.NewSelect().Model(&r).Where("url=?", url).Scan(context.Background())
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		return nil
	})
	if err != nil {
		return false, err
	}

	//検索したテーブルにデータが存在したのでfalseにする
	if r.Rss_URL != "" {
		return false, nil
	}

	//存在しなかったのでfalse
	return true, nil
}

// rssURLの登録
func registerRSS(r *tables.RSS_URLS) error {
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

//userItemに登録
func RegisterSubscriptionUserItem()error{
	//TODO:この続きから
	//ユーザーIDはクッキーからとれる
	//RSSIDはそのURLのが存在していたらそれを
	//していなかったらautoで作れるよう入力せず
	//でもgolangの場合自動で初期化され0になるのではという疑問
	//
	/*
	u := tables.USER_ITEMS{
		Rss_id: ,
		User_id: ,
	}
	err := database.WithDBConnection(func(db *bun.DB) error {
		_, err := db.NewInsert().Model(&r).Exec(context.Background())
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil{
		return err
	}

	return nil
	*/
}

