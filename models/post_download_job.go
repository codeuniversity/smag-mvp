package models

// PostDownloadJob represents a post which should be downloaded
type PostDownloadJob struct {
	PostID     int    `json:"post_id"`
	PictureURL string `json:"picture_url"`
}
