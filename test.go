package main

import (
	"rss"
	"flag"
	"fmt"
	"http"
	"time"
)

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Println("usage: test feed_url")
		return
	}
	r, err := http.DefaultClient.Get(flag.Arg(0))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	feed, err := rss.Get(r.Body)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Println(feed.Title)
	fmt.Println(feed.Subtitle)
	fmt.Println(feed.Link)
	for _, i := range feed.Items {
		fmt.Printf("\n")
		fmt.Println(i.Id)
		fmt.Println(i.Title)
		fmt.Println(i.Link)
		fmt.Println(time.SecondsToLocalTime(i.When).Format(time.UnixDate))
		fmt.Println(i.Description)
	}
}
