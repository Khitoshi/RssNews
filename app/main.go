package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"

	_ "github.com/lib/pq"
)

type Feed struct {
	Title           string //タイトル
	Link            string //記事リンク
	Description     string
	PublishedParsed *time.Time
	UpdatedParsed   *time.Time
}

// rss読み込み
func loadFeed() ([]Feed, error) {
	rssfeed, err := gofeed.NewParser().ParseURL("https://qiita.com/IXKGAGB/feed")
	if err != nil {
		//log.Fatal(err)
		return nil, err
	}

	//fmt.Println(feed.Title)
	fmt.Println(rssfeed.FeedType, rssfeed.FeedVersion)

	f := make([]Feed, len(rssfeed.Items))
	//feedに登録
	for i, item := range rssfeed.Items {

		f[i] = Feed{
			Title:           item.Title,
			Link:            item.Link,
			Description:     item.Description,
			PublishedParsed: item.PublishedParsed,
			UpdatedParsed:   item.UpdatedParsed,
		}

		fmt.Printf("title: %v\n", item.Title)
		fmt.Printf("\t-> %v\n", item.Link)
		fmt.Printf("\t-> %v\n", item.Description)     //説明
		fmt.Printf("\t-> %v\n", item.PublishedParsed) //記事アップ
		fmt.Printf("\t-> %v\n\n", item.UpdatedParsed) //最終更新
	}

	return f, nil
}

// templateにfeedを渡す
func MyHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/TestHome.html")
	if err != nil {
		panic(err)
	}

	feed, err := loadFeed()
	if err != nil {
		panic(err)
	}

	//HTMLテンプレートを実行
	err = tmpl.Execute(w, feed)
	if err != nil {
		panic(err)
	}
}

func updateFeed() error {
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

	rssfeed, err := gofeed.NewParser().ParseURL("https://qiita.com/IXKGAGB/feed")
	if err != nil {
		//log.Fatal(err)
		return err
	}
	for _, item := range rssfeed.Items {
		db.Exec("INSERT INTO items(url,title,description,author,published_at,created_at,updated_parsed) VALUES($1$2$3$4$5$6$7);",
			item.Link,
			item.Title,
			item.Description,
			item.Author,
			item.PublishedParsed,
			time.Now(),
			item.UpdatedParsed,
		)
	}

	return nil
}

func main() {
	startTime := time.Now()
	fmt.Printf("start time: %v \n", startTime)
	defer func() {
		fmt.Printf("\n processing time: %v", time.Since(startTime).Milliseconds())
	}()

	go updateFeed()

	http.HandleFunc("/", MyHandler)
	//ウェブサーバを起動
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
