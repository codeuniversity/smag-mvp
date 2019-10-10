package models

// TwitterPostRaw is the struct containing all raw twitter post fields
type TwitterPostRaw struct {
	// meta
	ID             int    `json:"id"`
	IDstr          string `json:"id_str"`
	ConversationID string `json:"conversation_id"`
	Link           string `json:"link"`
	Type           string `json:"tweet"`

	// time, place
	DateStamp string `json:"datestamp"`
	DateTime  int    `json:"datetime"`
	TimeStamp string `json:"timestamp"`
	TimeZone  string `json:"timezone"`
	Geo       string `json:"geo"`
	Near      string `json:"near"`
	Place     string `json:"place"`

	// content
	Cashtags    []string    `json:"cashtags`
	Hashtags    []string    `json:"hashtags`
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

// ReplyUser abc
type ReplyUser struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
}
