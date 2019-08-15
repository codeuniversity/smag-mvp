package main

import (
	"fmt"
	"strings"

	"github.com/alexmorten/instascraper/models"

	"github.com/gocolly/colly"
)

//ScrapeUserFollowGraph returns the follow information for a userName
func ScrapeUserFollowGraph(userName string) (*models.UserFollowInfo, error) {
	u := &models.UserFollowInfo{UserName: userName}

	followers, err := getUserNamesIn(fmt.Sprintf("http://picdeer.com/%s/followers", userName))
	if err != nil {
		return nil, err
	}
	u.Followers = followers
	followings, err := getUserNamesIn(fmt.Sprintf("http://picdeer.com/%s/followings", userName))
	if err != nil {
		return nil, err
	}
	u.Followings = followings
	return u, nil
}

func getUserNamesIn(url string) ([]string, error) {
	userNames := []string{}

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		// fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		// log.Println("Something went wrong:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		// fmt.Println("Visited", r.Request.URL)
	})

	c.OnHTML("p.grid-user-identifier-1 a", func(e *colly.HTMLElement) {
		url := e.Attr("href")
		parts := strings.Split(url, "/")
		if len(parts) >= 2 {
			name := parts[len(parts)-1]
			userNames = append(userNames, name)
		}
	})

	c.OnScraped(func(r *colly.Response) {
		// fmt.Println("Finished", r.Request.URL)
	})

	err := c.Visit(url)
	if err != nil {
		return nil, err
	}

	return userNames, nil
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
