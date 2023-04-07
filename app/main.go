package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"rss_reader/loadFeed"
	"rss_reader/updateFeed"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

type UserSignupParams struct {
	name      string `json:"name"`
	mail      string `json:"mail"`
	Apassword string `json:"apassword"`
	Bpassword string `json:"bpassword"`
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// 記事を表示するのページ
func home(c echo.Context) error {
	// 記事を取得
	feed, err := loadFeed.GetFeeds()
	if err != nil {
		panic(err)
	}
	return c.Render(http.StatusOK, "testhome", feed)
}

func veiwSignup(c echo.Context) error {
	return c.Render(http.StatusOK, "signup", nil)
}

// ユーザー登録するページ
func signup(c echo.Context) error {
	userparam := &UserSignupParams{}
	userparam.name = c.FormValue("name")
	userparam.mail = c.FormValue("mail")
	userparam.Apassword = c.FormValue("APassword")
	userparam.Bpassword = c.FormValue("BPassword")

	return c.Redirect(http.StatusFound, "signup")
}

func main() {
	startTime := time.Now()
	fmt.Printf("start time: %v \n", startTime)
	defer func() {
		fmt.Printf("\n processing time: %v", time.Since(startTime).Milliseconds())
	}()

	go updateFeed.UpdateItemsFromRSSFeed()

	e := echo.New()

	// ミドルウェアを設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Renderer = &Template{
		templates: template.Must(template.ParseGlob("./templates/*.html")),
	}

	e.GET("/home", home)

	e.GET("/signup", veiwSignup)
	e.POST("/signup", signup)

	// サーバーをポート番号8080で起動
	e.Logger.Fatal(e.Start(":8080"))
}
