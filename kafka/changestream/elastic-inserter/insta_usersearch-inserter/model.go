package main

// TODO: Update to user
type user struct {
	ID       int    `json:"id"`
	Username string `json:"user_name"`
	Realname string `json:"real_name"`
	Bio      string `json:"bio"`
}

// instaUsersearchUpsert is wrapped around the user (in json) to achieve upserting
// script - if user already exists
// upsert - if user doesn't exist yet
const instaUserUpsert = `
{
	"script" : {
        "source": "ctx._source.id = params.id; ctx._source.user_name = params.user_name; ctx._source.real_name = params.real_name",
        "lang": "painless",
        "params" : %s
    },
    "upsert" : %s
}
`

// TODO: Update to instaUsersearchMapping
const instaUserMapping = `
{
    "mappings" : {
		"properties" : {
			"id": {
				"type": "integer"
			}
			"user_name": {
				"type": "text"
			}
			"real_name": {
				"type": "text"
			}
			"bio": {
				"type": "text"
			}
		}
	}
}
`
