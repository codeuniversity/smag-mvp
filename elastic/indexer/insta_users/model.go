package main

type user struct {
	ID       int    `json:"id"`
	Username string `json:"user_name"`
	Realname string `json:"real_name"`
	Bio      string `json:"bio"`
}
