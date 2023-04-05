package main

import (
	"fmt"
	"html/template"
	"net/http"
	"rss_reader/loadFeed"
	"rss_reader/updateFeed"

	"time"

	_ "github.com/lib/pq"
)

func home(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/TestHome.html")
	if err != nil {
		panic(err)
	}

	//記事を取得
	feed, err := loadFeed.GetFeeds()
	if err != nil {
		panic(err)
	}

	//HTMLテンプレートを実行
	err = tmpl.Execute(w, feed)
	if err != nil {
		panic(err)
	}
}

func signup(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/Signup.html")
	if err != nil {
		panic(err)
	}

	//HTMLテンプレートを実行

	err = tmpl.Execute(w, nil)
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

	go updateFeed.UpdateItemsFromRSSFeed()

	//http.HandleFunc("/RSSReader", MyHandler)
	http.HandleFunc("/RSSReader/home", home)
	http.HandleFunc("/RSSReader/signup", signup)

	//ウェブサーバを起動
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}

}
