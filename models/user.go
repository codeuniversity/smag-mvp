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
