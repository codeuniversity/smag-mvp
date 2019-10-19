package models

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
	UID       string  `json:"uid,omitempty"`
	Name      string  `json:"name,omitempty"`
	RealName  string  `json:"real_name,omitempty"`
	AvatarURL string  `json:"avatar_url,omitempty"`
	Bio       string  `json:"bio,omitempty"`
	Follows   []*User `json:"follows,omitempty"`
	CrawledAt int     `json:"crawled_at,omitempty"`
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

// InstagramPost is a Post on instagram
type InstagramPost struct {
	PostID     string `json:"post_id"`
	ShortCode  string `json:"short_code"`
	UserID     string `json:"user_id"`
	PictureURL string `json:"picture_url"`
}

// InstaCommentScrapError s are written to the error topic of the scraper when even after retries we can't scrape the comments
type InstaCommentScrapError struct {
	PostID string `json:"post_id,omitempty"`
	Error  string `json:"error,omitempty"`
}
