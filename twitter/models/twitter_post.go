package models

import (
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"
)

// TwitterPostRaw is the struct containing all raw twitter post fields
type TwitterPostRaw struct {
	// meta
	ID             int    `json:"id"`
	IDstr          string `json:"id_str"`
	ConversationID string `json:"conversation_id"`
	Link           string `json:"link"`
	Type           string `json:"type"`

	// time, place
	DateStamp string `json:"datestamp"`
	DateTime  int    `json:"datetime"`
	TimeStamp string `json:"timestamp"`
	TimeZone  string `json:"timezone"`
	Geo       string `json:"geo"`
	Near      string `json:"near"`
	Place     string `json:"place"`

	// content
	Cashtags    []string    `json:"cashtags"`
	Hashtags    []string    `json:"hashtags"`
	Mentions    []string    `json:"mentions"`
	Photos      []string    `json:"photos"`
	QuoteURL    string      `json:"quote_url"`
	ReplyTo     []ReplyUser `json:"reply_to"`
	Retweet     bool        `json:"retweet"`
	RetweetDate string      `json:"retweet_date"`
	RetweetID   string      `json:"retweet_id"`
	Source      string      `json:"source"`
	Tweet       string      `json:"tweet"`
	URLs        []string    `json:"urls"`
	Video       int         `json:"video"`

	// reactions
	LikesCount   string `json:"likes_count"`
	RepliesCount string `json:"replies_count"`
	RetweetCount string `json:"retweet_count"`

	// user info
	Name      string `json:"name"`
	UserID    int    `json:"user_id"`
	UserIDstr string `json:"user_id_str"`
	UserName  string `json:"username"`
	UserRt    string `json:"user_rt"`
	UserRtID  string `json:"user_tr_id"`
}

// TwitterPost is the struct containing all processed twitter post fields
type TwitterPost struct {
	GormModelWithoutID

	// meta
	PostIdentifier uint64 `gorm:"primary_key:true"`
	ConversationID string
	Link           string
	Type           string

	// time, place
	DateTime time.Time
	TimeZone string
	Geo      string
	Near     string
	Place    string

	// content
	Cashtags    pq.StringArray `gorm:"type:text[]"`
	Hashtags    pq.StringArray `gorm:"type:text[]"`
	Mentions    []*TwitterUser `gorm:"many2many:post_mentions;association_foreignkey:username;foreignkey:postIdentifier"`
	Photos      pq.StringArray `gorm:"type:text[]"`
	QuoteURL    string
	ReplyTo     []*TwitterUser `gorm:"many2many:post_replies;association_foreignkey:username;foreignkey:postIdentifier"`
	Retweet     bool
	RetweetDate string
	RetweetID   string
	Source      string
	Tweet       string
	URLs        pq.StringArray `gorm:"type:text[]"`
	Video       int

	// reactions
	LikesCount   int
	RepliesCount int
	RetweetCount int

	UserID          string
	Username        string
	RetweetUserID   string
	RetweetUsername string
}

// ReplyUser abc
type ReplyUser struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
}

// ConvertTwitterPost converts the raw TwitterPost structure
// from kafka into the database model
func ConvertTwitterPost(raw *TwitterPostRaw) *TwitterPost {
	mentions := make([]*TwitterUser, len(raw.Mentions))
	replyTo := make([]*TwitterUser, len(raw.ReplyTo))

	for index, username := range raw.Mentions {
		mentions[index] = &TwitterUser{
			Username: strings.ToLower(username),
		}
	}

	for index, replyUser := range raw.ReplyTo {
		replyTo[index] = &TwitterUser{
			Username: strings.ToLower(replyUser.Username),
		}
	}

	dateTime := time.Unix(int64(raw.DateTime/1000), 0)

	likesCount, _ := strconv.Atoi(raw.LikesCount)
	repliesCount, _ := strconv.Atoi(raw.RepliesCount)
	retweetCount, _ := strconv.Atoi(raw.RetweetCount)

	return &TwitterPost{
		PostIdentifier: uint64(raw.ID),
		ConversationID: raw.ConversationID,
		Link:           raw.Link,
		Type:           raw.Type,

		DateTime: dateTime,
		TimeZone: raw.TimeZone,
		Geo:      raw.Geo,
		Near:     raw.Near,
		Place:    raw.Place,

		Cashtags:    raw.Cashtags,
		Hashtags:    raw.Hashtags,
		Mentions:    mentions,
		Photos:      raw.Photos,
		QuoteURL:    raw.QuoteURL,
		ReplyTo:     replyTo,
		Retweet:     raw.Retweet,
		RetweetDate: raw.RetweetDate,
		RetweetID:   raw.RetweetID,
		Source:      raw.Source,
		Tweet:       raw.Tweet,
		URLs:        raw.URLs,
		Video:       raw.Video,

		LikesCount:   likesCount,
		RepliesCount: repliesCount,
		RetweetCount: retweetCount,

		UserID:          raw.UserIDstr,
		Username:        strings.ToLower(raw.UserName),
		RetweetUserID:   raw.UserRtID,
		RetweetUsername: strings.ToLower(raw.UserRt),
	}
}
