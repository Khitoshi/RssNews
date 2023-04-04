package main

import (
	"fmt"
	"html/template"
	"net/http"
	item "rss_reader/updateFeed"
	"time"

	_ "github.com/lib/pq"
	"github.com/mmcdole/gofeed"
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

func main() {
	startTime := time.Now()
	fmt.Printf("start time: %v \n", startTime)
	defer func() {
		fmt.Printf("\n processing time: %v", time.Since(startTime).Milliseconds())
	}()

	go item.FixedTermUpdate()

	http.HandleFunc("/", MyHandler)
	//ウェブサーバを起動
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
