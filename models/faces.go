package models

import (
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
)

// FaceData represents the face encoding table
type FaceData struct {
	gorm.Model
	PostID   int            `json:"post_id"`
	X        int            `json:"x"`
	Y        int            `json:"y"`
	Width    int            `json:"width"`
	Height   int            `json:"height"`
	Encoding postgres.Jsonb `json:"encoding"`
}

// FaceRecognitionResult is the result of the face recognizer
type FaceRecognitionResult struct {
	PostID int     `json:"post_id"`
	Faces  []*Face `json:"faces"`
}

// FaceReconJob represents the data for a face recon job
type FaceReconJob struct {
	PostID           int    `json:"post_id"`
	InternalImageURL string `json:"internal_image_url"`
	X                int    `json:"x"`
	Y                int    `json:"y"`
	Width            int    `json:"width"`
	Height           int    `json:"height"`
}

// Face contains the position of a face in a post and its encoding
type Face struct {
	X        int          `json:"x"`
	Y        int          `json:"y"`
	Width    int          `json:"width"`
	Height   int          `json:"height"`
	Encoding [128]float32 `json:"encoding"`
}
