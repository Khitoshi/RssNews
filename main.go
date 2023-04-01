package main

import (
	"fmt"
	"log"
	"time"

	"github.com/mmcdole/gofeed"
)

func main() {
	startTime := time.Now()
	fmt.Printf("start time: %v \n", startTime)
	defer func() {
		fmt.Printf("\n processing time: %v", time.Since(startTime).Milliseconds())
	}()

	feed, err := gofeed.NewParser().ParseURL("https://qiita.com/IXKGAGB/feed")
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println(feed.Title)
	fmt.Println(feed.FeedType, feed.FeedVersion)

	for _, item := range feed.Items {
		fmt.Printf("title: %v\n", item.Title)
		fmt.Printf("\t-> %v\n", item.Link)
		fmt.Printf("\t-> %v\n", item.Description)
		//fmt.Println(item.PubDate)
		fmt.Printf("\t-> %v\n", item.PublishedParsed)
		fmt.Printf("\t-> %v\n\n", item.UpdatedParsed)
		//fmt.Println(item.CreatedParsed)
	}

}
