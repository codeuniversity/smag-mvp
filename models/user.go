package models

//UserFollowInfo holds the follow graph info, only relating userNames
type UserFollowInfo struct {
	UserName   string   `json:"user_name"`
	Followers  []string `json:"followers"`
	Followings []string `json:"followings"`
}

// User is the struct containing all user fields, used for serializing users to dgraph
type User struct {
	UID        string  `json:"uid,omitempty"`
	Name       string  `json:"name,omitempty"`
	Followers  []*User `json:"follows,omitempty"`
	Followings []*User `json:"followed,omitempty"`
}
