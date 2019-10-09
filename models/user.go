package models

import "github.com/jinzhu/gorm"

//UserFollowInfo holds the follow graph info, only relating userNames
type UserFollowInfo struct {
	gorm.Model
	UserName   string   `json:"user_name" gorm:"column:user_name"`
	RealName   string   `json:"real_name" gorm:"column:real_name"`
	AvatarURL  string   `json:"avatar_url" gorm:"column:avatar_url"`
	Bio        string   `json:"bio" gorm:"column:bio"`
	Followers  []string `json:"followers" gorm:"column:followers"`
	Followings []string `json:"followings" gorm:"column:followings"`
	CrawlTs    int      `json:"crawl_ts" gorm:"column:crawl_ts"`
}

// User is the struct containing all user fields, used for serializing users to dgraph
type User struct {
	gorm.Model
	UID       string  `json:"uid,omitempty"`
	UserName  string  `json:"user_name,omitempty" gorm:"column:user_name"`
	RealName  string  `json:"real_name,omitempty" gorm:"column:real_name"`
	AvatarURL string  `json:"avatar_url,omitempty" gorm:"column:avatar_url"`
	Bio       string  `json:"bio,omitempty" gorm:"column:bio"`
	Follows   []*User `json:"follows,omitempty" gorm:"column:follows"`
	CrawledAt int     `json:"crawled_at,omitempty" gorm:"column:crawl_ts"`
}

// Follow is a struct representing the "follows" table in the postgres database
type Follow struct {
	// UID  string `json:"uid,omitempty" gorm:"column:uid"`
	From uint `gorm:"column:from_id"`
	To   uint `gorm:"column:to_id"`
}

// ScrapeError s are written to user_scrape_errors when even after retries we can't scrape a user
type ScrapeError struct {
	Name  string `json:"name,omitempty"`
	Error string `json:"error,omitempty"`
}
