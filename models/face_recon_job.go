package models

// FaceReconJob represents the data for a face recon job
type FaceReconJob struct {
	PostID           int    `json:"post_id"`
	InternalImageURL string `json:"internal_image_url"`
	X                int    `json:"x"`
	Y                int    `json:"y"`
	Width            int    `json:"width"`
	Height           int    `json:"height"`
}
