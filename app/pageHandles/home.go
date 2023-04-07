package pageHandles

import (
	"net/http"
	"rss_reader/loadFeed"

	"github.com/labstack/echo/v4"
)

// 記事を表示するのページ
func HandleHome_Get(c echo.Context) error {
	// 記事を取得
	feed, err := loadFeed.GetFeeds()
	if err != nil {
		panic(err)
	}
	return c.Render(http.StatusOK, "testhome", feed)
}
