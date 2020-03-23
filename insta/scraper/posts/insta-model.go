package scraper

type instagramAccountInfo struct {
	LoggingPageID         string `json:"logging_page_id"`
	ShowSuggestedProfiles bool   `json:"show_suggested_profiles"`
	ShowFollowDialog      bool   `json:"show_follow_dialog"`
	Graphql               struct {
		User struct {
			Biography              string `json:"biography"`
			BlockedByViewer        bool   `json:"blocked_by_viewer"`
			CountryBlock           bool   `json:"country_block"`
			ExternalURL            string `json:"external_url"`
			ExternalURLLinkshimmed string `json:"external_url_linkshimmed"`
			EdgeFollowedBy         struct {
				Count int `json:"count"`
			} `json:"edge_followed_by"`
			FollowedByViewer bool `json:"followed_by_viewer"`
			EdgeFollow       struct {
				Count int `json:"count"`
			} `json:"edge_follow"`
			FollowsViewer        bool   `json:"follows_viewer"`
			FullName             string `json:"full_name"`
			HasChannel           bool   `json:"has_channel"`
			HasBlockedViewer     bool   `json:"has_blocked_viewer"`
			HighlightReelCount   int    `json:"highlight_reel_count"`
			HasRequestedViewer   bool   `json:"has_requested_viewer"`
			ID                   string `json:"id"`
			IsBusinessAccount    bool   `json:"is_business_account"`
			IsJoinedRecently     bool   `json:"is_joined_recently"`
			BusinessCategoryName string `json:"business_category_name"`
			IsPrivate            bool   `json:"is_private"`
			IsVerified           bool   `json:"is_verified"`
			EdgeMutualFollowedBy struct {
				Count int           `json:"count"`
				Edges []interface{} `json:"edges"`
			} `json:"edge_mutual_followed_by"`
			ProfilePicURL          string      `json:"profile_pic_url"`
			ProfilePicURLHd        string      `json:"profile_pic_url_hd"`
			RequestedByViewer      bool        `json:"requested_by_viewer"`
			Username               string      `json:"username"`
			ConnectedFbPage        interface{} `json:"connected_fb_page"`
			EdgeFelixVideoTimeline struct {
				Count    int `json:"count"`
				PageInfo struct {
					HasNextPage bool        `json:"has_next_page"`
					EndCursor   interface{} `json:"end_cursor"`
				} `json:"page_info"`
				Edges []interface{} `json:"edges"`
			} `json:"edge_felix_video_timeline"`
			EdgeOwnerToTimelineMedia struct {
				Count    int `json:"count"`
				PageInfo struct {
					HasNextPage bool   `json:"has_next_page"`
					EndCursor   string `json:"end_cursor"`
				} `json:"page_info"`
				Edges []struct {
					Node struct {
						Typename           string `json:"__typename"`
						ID                 string `json:"id"`
						EdgeMediaToCaption struct {
							Edges []struct {
								Node struct {
									Text string `json:"text"`
								} `json:"node"`
							} `json:"edges"`
						} `json:"edge_media_to_caption"`
						Shortcode          string `json:"shortcode"`
						EdgeMediaToComment struct {
							Count int `json:"count"`
						} `json:"edge_media_to_comment"`
						CommentsDisabled bool `json:"comments_disabled"`
						TakenAtTimestamp int  `json:"taken_at_timestamp"`
						Dimensions       struct {
							Height int `json:"height"`
							Width  int `json:"width"`
						} `json:"dimensions"`
						DisplayURL  string `json:"display_url"`
						EdgeLikedBy struct {
							Count int `json:"count"`
						} `json:"edge_liked_by"`
						EdgeMediaPreviewLike struct {
							Count int `json:"count"`
						} `json:"edge_media_preview_like"`
						Location struct {
							ID            string `json:"id"`
							HasPublicPage bool   `json:"has_public_page"`
							Name          string `json:"name"`
							Slug          string `json:"slug"`
						} `json:"location"`
						GatingInfo           interface{} `json:"gating_info"`
						FactCheckInformation interface{} `json:"fact_check_information"`
						MediaPreview         interface{} `json:"media_preview"`
						Owner                struct {
							ID       string `json:"id"`
							Username string `json:"username"`
						} `json:"owner"`
						ThumbnailSrc       string `json:"thumbnail_src"`
						ThumbnailResources []struct {
							Src          string `json:"src"`
							ConfigWidth  int    `json:"config_width"`
							ConfigHeight int    `json:"config_height"`
						} `json:"thumbnail_resources"`
						IsVideo              bool   `json:"is_video"`
						AccessibilityCaption string `json:"accessibility_caption"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"edge_owner_to_timeline_media"`
			EdgeSavedMedia struct {
				Count    int `json:"count"`
				PageInfo struct {
					HasNextPage bool        `json:"has_next_page"`
					EndCursor   interface{} `json:"end_cursor"`
				} `json:"page_info"`
				Edges []interface{} `json:"edges"`
			} `json:"edge_saved_media"`
			EdgeMediaCollections struct {
				Count    int `json:"count"`
				PageInfo struct {
					HasNextPage bool        `json:"has_next_page"`
					EndCursor   interface{} `json:"end_cursor"`
				} `json:"page_info"`
				Edges []interface{} `json:"edges"`
			} `json:"edge_media_collections"`
		} `json:"user"`
	} `json:"graphql"`
	ToastContentOnLoad interface{} `json:"toast_content_on_load"`
}

type instagramMedia struct {
	Data struct {
		User struct {
			EdgeOwnerToTimelineMedia struct {
				Count    int `json:"count"`
				PageInfo struct {
					HasNextPage bool   `json:"has_next_page"`
					EndCursor   string `json:"end_cursor"`
				} `json:"page_info"`
				Edges []struct {
					Node struct {
						Typename   string `json:"__typename"`
						ID         string `json:"id"`
						Dimensions struct {
							Height int `json:"height"`
							Width  int `json:"width"`
						} `json:"dimensions"`
						DisplayURL       string `json:"display_url"`
						DisplayResources []struct {
							Src          string `json:"src"`
							ConfigWidth  int    `json:"config_width"`
							ConfigHeight int    `json:"config_height"`
						} `json:"display_resources"`
						IsVideo               bool   `json:"is_video"`
						TrackingToken         string `json:"tracking_token"`
						EdgeMediaToTaggedUser struct {
							Edges []struct {
								Node struct {
									User struct {
										FullName      string `json:"full_name"`
										ID            string `json:"id"`
										IsVerified    bool   `json:"is_verified"`
										ProfilePicURL string `json:"profile_pic_url"`
										Username      string `json:"username"`
									} `json:"user"`
									X float64 `json:"x"`
									Y float64 `json:"y"`
								} `json:"node"`
							} `json:"edges"`
						} `json:"edge_media_to_tagged_user"`
						AccessibilityCaption interface{} `json:"accessibility_caption"`
						EdgeMediaToCaption   struct {
							Edges []struct {
								Node struct {
									Text string `json:"text"`
								} `json:"node"`
							} `json:"edges"`
						} `json:"edge_media_to_caption"`
						Shortcode          string `json:"shortcode"`
						EdgeMediaToComment struct {
							Count    int `json:"count"`
							PageInfo struct {
								HasNextPage bool        `json:"has_next_page"`
								EndCursor   interface{} `json:"end_cursor"`
							} `json:"page_info"`
							Edges []struct {
								Node struct {
									ID              string `json:"id"`
									Text            string `json:"text"`
									CreatedAt       int    `json:"created_at"`
									DidReportAsSpam bool   `json:"did_report_as_spam"`
									Owner           struct {
										ID            string `json:"id"`
										IsVerified    bool   `json:"is_verified"`
										ProfilePicURL string `json:"profile_pic_url"`
										Username      string `json:"username"`
									} `json:"owner"`
									ViewerHasLiked bool `json:"viewer_has_liked"`
								} `json:"node"`
							} `json:"edges"`
						} `json:"edge_media_to_comment"`
						EdgeMediaToSponsorUser struct {
							Edges []interface{} `json:"edges"`
						} `json:"edge_media_to_sponsor_user"`
						CommentsDisabled     bool `json:"comments_disabled"`
						TakenAtTimestamp     int  `json:"taken_at_timestamp"`
						EdgeMediaPreviewLike struct {
							Count int           `json:"count"`
							Edges []interface{} `json:"edges"`
						} `json:"edge_media_preview_like"`
						GatingInfo           interface{} `json:"gating_info"`
						FactCheckInformation interface{} `json:"fact_check_information"`
						MediaPreview         interface{} `json:"media_preview"`
						Owner                struct {
							ID       string `json:"id"`
							Username string `json:"username"`
						} `json:"owner"`
						Location struct {
							ID            string `json:"id"`
							HasPublicPage bool   `json:"has_public_page"`
							Name          string `json:"name"`
							Slug          string `json:"slug"`
						} `json:"location"`
						ViewerHasLiked             bool   `json:"viewer_has_liked"`
						ViewerHasSaved             bool   `json:"viewer_has_saved"`
						ViewerHasSavedToCollection bool   `json:"viewer_has_saved_to_collection"`
						ViewerInPhotoOfYou         bool   `json:"viewer_in_photo_of_you"`
						ViewerCanReshare           bool   `json:"viewer_can_reshare"`
						ThumbnailSrc               string `json:"thumbnail_src"`
						ThumbnailResources         []struct {
							Src          string `json:"src"`
							ConfigWidth  int    `json:"config_width"`
							ConfigHeight int    `json:"config_height"`
						} `json:"thumbnail_resources"`
						EdgeSidecarToChildren struct {
							Edges []struct {
								GraphVideo struct {
									Typename   string `json:"__typename"`
									ID         string `json:"id"`
									Dimensions struct {
										Height int `json:"height"`
										Width  int `json:"width"`
									} `json:"dimensions"`
									DisplayURL       string `json:"display_url"`
									DisplayResources []struct {
										Src          string `json:"src"`
										ConfigWidth  int    `json:"config_width"`
										ConfigHeight int    `json:"config_height"`
									} `json:"display_resources"`
									IsVideo               bool   `json:"is_video"`
									TrackingToken         string `json:"tracking_token"`
									EdgeMediaToTaggedUser struct {
										Edges []interface{} `json:"edges"`
									} `json:"edge_media_to_tagged_user"`
									DashInfo struct {
										IsDashEligible    bool        `json:"is_dash_eligible"`
										VideoDashManifest interface{} `json:"video_dash_manifest"`
										NumberOfQualities int         `json:"number_of_qualities"`
									} `json:"dash_info"`
									VideoURL       string `json:"video_url"`
									VideoViewCount int    `json:"video_view_count"`
								} `json:"node,omitempty"`
							} `json:"edges"`
						} `json:"edge_sidecar_to_children"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"edge_owner_to_timeline_media"`
		} `json:"user"`
	} `json:"data"`
	Status string `json:"status"`
}
