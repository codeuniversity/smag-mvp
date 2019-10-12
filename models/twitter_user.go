package models

// TwitterUserRaw holds the follow graph info, only relating userNames
type TwitterUserRaw struct {
	// Meta
	ID   string `json:"id"`
	URL  string `json:"url"`
	Type string `json:"type"`

	// User info
	Name            string `json:"name"`
	Username        string `json:"username"`
	Bio             string `json:"bio"`
	Avatar          string `json:"avatar"`
	BackgroundImage string `json:"background_image"`

	// Profile stats
	Location   string `json:"location"`
	JoinDate   string `json:"join_date"`
	JoinTime   string `json:"join_time"`
	IsPrivate  int    `json:"is_private"`
	IsVerified int    `json:"is_verified"`

	// Follows
	Following     int      `json:"following"`
	FollowingList []string `json:"following_list"`
	Followers     int      `json:"followers"`
	FollowersList []string `json:"followers_list"`

	// Usage stats
	Tweets     int `json:"tweets"`
	Likes      int `json:"likes"`
	MediaCount int `json:"media_count"`
}

// TwitterUser holds the follow graph info, only relating userNames
type TwitterUser struct {
	GormModelWithoutID

	// Meta
	ID   string
	URL  string
	Type string

	// User info
	Name            string
	Username        string
	Bio             string
	Avatar          string
	BackgroundImage string `gorm:"column:background_image"`

	// Profile stats
	Location   string
	JoinDate   string `json:"join_date"`
	JoinTime   string `json:"join_time"`
	IsPrivate  int    `json:"is_private"`
	IsVerified int    `json:"is_verified"`

	// Follows
	Following     int
	FollowingList []string `json:"following_list"`
	Followers     int
	FollowersList []string `json:"followers_list"`

	// Usage stats
	Tweets     int
	Likes      int `json:"likes"`
	MediaCount int `json:"media_count"`
}
