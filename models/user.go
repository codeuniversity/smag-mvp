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

// User is the struct containing all user fields, used for serializing users to dgraph
type User struct {
	UID       string  `json:"uid,omitempty"`
	Name      string  `json:"name,omitempty"`
	RealName  string  `json:"real_name,omitempty"`
	AvatarURL string  `json:"avatar_url,omitempty"`
	Bio       string  `json:"bio,omitempty"`
	Follows   []*User `json:"follows,omitempty"`
	CrawledAt int     `json:"crawled_at,omitempty"`
}

type RenewingAddresses struct {
	InstanceId string   `json:"instanceId"`
	LocalIps   []string `json:"localIps"`
}

// ScrapeError s are written to user_scrape_errors when even after retries we can't scrape a user
type ScrapeError struct {
	Name  string `json:"name,omitempty"`
	Error string `json:"error,omitempty"`
}

// ScrapeError s are written to user_scrape_errors when even after retries we can't scrape a user
type AwsServiceError struct {
	InstanceId string `json:"name,omitempty"`
	Error      string `json:"error,omitempty"`
}

// ScrapeError s are written to user_scrape_errors when even after retries we can't scrape a user
type InstagramScrapeError struct {
	Name  string `json:"name,omitempty"`
	Error string `json:"error,omitempty"`
}

type InstaComment struct {
	Id            string `json:"id"`
	Text          string `json:"text"`
	CreatedAt     int    `json:"created_at"`
	PostId        string `json:"post_id"`
	ShortCode     string `json:"short_code"`
	OwnerUsername string `json:"owner_username"`
}

type InstagramPost struct {
	PostId     string `json:"post_id"`
	ShortCode  string `json:"short_code"`
	UserId     string `json:"user_id"`
	PictureUrl string `json:"picture_url"`
}

type InstaCommentScrapError struct {
	PostId string `json:"post_id,omitempty"`
	Error  string `json:"error,omitempty"`
}
