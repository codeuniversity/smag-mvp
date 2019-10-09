package main

import (
	"fmt"
	insta_post_comment_scraper "github.com/codeuniversity/smag-mvp/insta-post-comment-scraper"
)

func main() {

	client := insta_post_comment_scraper.New("52.58.171.160:9092")
	comments, err := client.Scrape("B2hByt_obIu")

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(comments.Status)
}
