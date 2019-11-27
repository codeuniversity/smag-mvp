package models

type InstaUser struct {
	ID       int    `json:"id"`
	Username string `json:"user_name"`
	Realname string `json:"real_name"`
	Bio      string `json:"bio"`
}

type InstaPost struct {
	ID      int    `json:"id"`
	UserID  string `json:"user_id"`
	Caption string `json:"caption"`
}

type InstaComment struct {
	ID      int    `json:"id"`
	PostID  string `json:"post_id"`
	Comment string `json:"comment_text"`
}
