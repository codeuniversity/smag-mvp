package insta_likes_scraper

type InstaPostLikes struct {
	Data struct {
		ShortcodeMedia struct {
			ID          string `json:"id"`
			Shortcode   string `json:"shortcode"`
			EdgeLikedBy struct {
				Count    int `json:"count"`
				PageInfo struct {
					HasNextPage bool   `json:"has_next_page"`
					EndCursor   string `json:"end_cursor"`
				} `json:"page_info"`
				Edges []struct {
					Node struct {
						ID                string `json:"id"`
						Username          string `json:"username"`
						FullName          string `json:"full_name"`
						ProfilePicURL     string `json:"profile_pic_url"`
						IsPrivate         bool   `json:"is_private"`
						IsVerified        bool   `json:"is_verified"`
						FollowedByViewer  bool   `json:"followed_by_viewer"`
						RequestedByViewer bool   `json:"requested_by_viewer"`
						Reel              struct {
							ID              string      `json:"id"`
							ExpiringAt      int         `json:"expiring_at"`
							HasPrideMedia   bool        `json:"has_pride_media"`
							LatestReelMedia interface{} `json:"latest_reel_media"`
							Seen            interface{} `json:"seen"`
							Owner           struct {
								Typename      string `json:"__typename"`
								ID            string `json:"id"`
								ProfilePicURL string `json:"profile_pic_url"`
								Username      string `json:"username"`
							} `json:"owner"`
						} `json:"reel"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"edge_liked_by"`
		} `json:"shortcode_media"`
	} `json:"data"`
	Status string `json:"status"`
}
