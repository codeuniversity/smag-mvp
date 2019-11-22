package insta_post_text_filter

type QueryResult struct {
	Took int `json:"took"`
	Hits struct {
		Total struct {
			Value    int    `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		MaxScore float64 `json:"max_score"`
		Hits     []struct {
			Index  string  `json:"_index"`
			Type   string  `json:"_type"`
			ID     string  `json:"_id"`
			Score  float64 `json:"_score"`
			Source struct {
				PostID string `json:"postId"`
				Text   string `json:"text"`
			} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

type InstaPost struct {
	PostId  string `json:"post_id"`
	Caption string `json:"caption"`
}
