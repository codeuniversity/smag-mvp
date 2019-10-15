package models

import (
	"time"

	"github.com/jinzhu/gorm"

	"github.com/codeuniversity/smag-mvp/utils"
)

// TwitterUserList is a custom type of TwitterUser to be used for easier handling
// of relation users in twitter inserters
type TwitterUserList []*TwitterUser

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
	gorm.Model

	// Meta
	TwitterID string
	URL       string
	Type      string

	// User info
	Name            string
	Username        string
	Bio             string
	Avatar          string
	BackgroundImage string

	// Profile stats
	Location   string
	JoinDate   time.Time
	IsPrivate  bool
	IsVerified bool

	// Follows
	Following     int
	FollowingList []*TwitterUser `gorm:"many2many:twitter_followings;"`
	Followers     int
	FollowersList []*TwitterUser `gorm:"many2many:twitter_followers;"`

	// Usage stats
	Tweets     int
	Likes      int
	MediaCount int
}

// ConvertTwitterUser converts the raw TwitterUser structure
// from kafka into the database model
func ConvertTwitterUser(raw *TwitterUserRaw) *TwitterUser {

	followingList := make([]*TwitterUser, len(raw.FollowingList))
	followersList := make([]*TwitterUser, len(raw.FollowersList))

	for index, item := range raw.FollowingList {
		followingList[index] = &TwitterUser{
			Username: item,
		}
	}

	for index, item := range raw.FollowingList {
		followingList[index] = &TwitterUser{
			Username: item,
		}
	}

	joinDate, _ := utils.ConvertDateStrToTime(raw.JoinDate)

	isPrivate := utils.ConvertIntToBool(raw.IsPrivate)
	isVerified := utils.ConvertIntToBool(raw.IsVerified)

	return &TwitterUser{
		TwitterID: raw.ID,
		URL:       raw.URL,
		Type:      raw.Type,

		Name:            raw.Name,
		Username:        raw.Username,
		Bio:             raw.Bio,
		Avatar:          raw.Avatar,
		BackgroundImage: raw.BackgroundImage,

		Location:   raw.Location,
		JoinDate:   joinDate,
		IsPrivate:  isPrivate,
		IsVerified: isVerified,

		Following:     raw.Following,
		FollowingList: followingList,
		Followers:     raw.Followers,
		FollowersList: followersList,

		Tweets:     raw.Tweets,
		Likes:      raw.Likes,
		MediaCount: raw.MediaCount,
	}
}

func (list *TwitterUserList) Create(slices ...[]*TwitterUser) {
	for _, slice := range slices {
		*list = append(*list, slice...)
	}
}

func (list *TwitterUserList) RemoveDuplicates() {
	set := make(map[*TwitterUser]bool)
	for index := range *list {
		if !set[(*list)[index]] {
			set[(*list)[index]] = true
		}
	}

	uniqueList := make([]*TwitterUser, len(set))
	index := 0
	for user := range set {
		uniqueList[index] = user
		index++
	}

	*list = uniqueList
}
