package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"time"

	"rss_reader/pageHandles"
	"rss_reader/tables"
	"rss_reader/updateFeed"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	//スタート時間・処理時間表示
	startTime := time.Now()
	fmt.Printf("start time: %v \n", startTime)
	defer func() {
		fmt.Printf("\n processing time: %v", time.Since(startTime).Milliseconds())
	}()

	//テーブル作成
	err := tables.CreateTable()
	if err != nil {
		log.Fatal(err)
	}

	//記事アップデート
	go updateFeed.RegularUpdatingOfArticles()

	e := echo.New()

	// ミドルウェアを設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/tmp", "./templates/")

	e.Renderer = &Template{
		templates: template.Must(template.ParseGlob("./templates/*.html")),
	}

	e.GET("/", pageHandles.Hoge)

	//e.GET("/", pageHandles.HandleHome_Get)
	e.GET("/home", pageHandles.HandleHome_Get)

	e.GET("/signup", pageHandles.HandleSignup_Get)
	e.POST("/signup", pageHandles.HandleSignup_Post)

	e.GET("/login", pageHandles.HandleLogin_Get)
	e.POST("/login", pageHandles.HandleLogin_Post)

	//ログアウト
	e.GET("/logout", pageHandles.HandleLogout_Get)
	e.POST("/logout", pageHandles.HandleLogout_Post)

	e.GET("/registerRSS", pageHandles.HandleRegisterRSS_Get)
	e.POST("/registerRSS", pageHandles.HandleRegisterRSS_Post)

	// サーバーをポート番号8080で起動
	e.Logger.Fatal(e.Start(":8080"))
}
