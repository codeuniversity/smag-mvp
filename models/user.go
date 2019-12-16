package models

import "github.com/jinzhu/gorm"

//UserFollowInfo holds the follow graph info, only relating userNames
type UserFollowInfo struct {
	UserName   string   `json:"user_name"`
	RealName   string   `json:"real_name"`
	AvatarURL  string   `json:"avatar_url"`
	Bio        string   `json:"bio"`
	Followers  []string `json:"followers"`
	Followings []string `json:"followings"`
	CrawlTs    int      `json:"crawl_ts"`
}

// User is the struct containing all user fields, used for saving users to the database
type User struct {
	gorm.Model
	UserName  string `gorm:"unique"`
	RealName  string
	AvatarURL string `gorm:"column:avatar_url"`
	Bio       string
	Follows   []*User
	CrawledAt int `gorm:"column:crawl_ts"`
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

// InstagramScrapeError s are written to user_scrape_errors when even after retries we can't scrape a user
type InstagramScrapeError struct {
	Name  string `json:"name,omitempty"`
	Error string `json:"error,omitempty"`
}

// InstaComment os a comment on instagram
type InstaComment struct {
	ID            string `json:"id"`
	Text          string `json:"text"`
	CreatedAt     int    `json:"created_at"`
	PostID        string `json:"post_id"`
	ShortCode     string `json:"short_code"`
	OwnerUsername string `json:"owner_username"`
}

// Comment is the database model used for creating and updating comments in postgres
type Comment struct {
	gorm.Model
	PostID      int
	CommentID   string
	CommentText string `gorm:"type:text"`
	OwnerUserID int    `gorm:"column:owner_user_id"`
}

type InstaLike struct {
	ID            string `json:"id"`
	PostID        string `json:"post_id"`
	OwnerUsername string `json:"owner_username"`
	AvatarURL     string `json:"avatar_url"`
	RealName      string `json:"real_name"`
}

// InstagramPost is a Post on instagram
type InstagramPost struct {
	gorm.Model
	PostID      string   `json:"post_id"`
	ShortCode   string   `json:"short_code"`
	UserName    string   `json:"user_name"`
	UserID      string   `json:"user_id"`
	PictureURL  string   `json:"picture_url"`
	TaggedUsers []string `json:"tagged_users"`
	Caption     string   `json:"caption"`
}

// Post is the database model used for creating and updating posts in postgres
type Post struct {
	gorm.Model
	UserID             int
	PostID             string
	ShortCode          string
	PictureURL         string `gorm:"column:picture_url;type:text"`
	Caption            string `gorm:"type:text"`
	InternalPictureURL string `gorm:"column:internal_picture_url"`
}

// InstaPostScrapeError s are written to the error topic of the scraper when even after retries we can't scrape the comments
type InstaPostScrapeError struct {
	PostID string `json:"post_id,omitempty"`
	Error  string `json:"error,omitempty"`
}

type AwsServiceError struct {
	InstanceId string `json:"instance_id"`
	Error      string `json:"error"`
}
