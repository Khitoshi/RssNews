package pageHandles

import (
	"net/http"
	"rss_reader/loadFeed"
	"strconv"

	"github.com/labstack/echo/v4"
)

// 記事を表示するのページ
func HandleHome_Get(c echo.Context) error {
	//クッキーからuseridを取得
	cookie, err := c.Cookie("userId")
	if err != nil {
		return c.Redirect(http.StatusFound, "login")
	}

	//クッキーから入手したユーザーidをint64に変換
	userId, err := strconv.ParseInt(cookie.Value, 10, 64)
	if err != nil {
		return err
	}

	// 記事を取得
	news, err := loadFeed.GetFeeds(userId)

	//feed, err := loadFeed.GetFeeds(8)
	if err != nil {
		panic(err)
	}
	return c.Render(http.StatusOK, "index", news)
}
