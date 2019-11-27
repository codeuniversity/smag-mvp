package models

// FaceDoc is the type that is used to store a face in elasticsearch
type FaceDoc struct {
	PostID         int    `json:"post_id"`
	X              int    `json:"x"`
	Y              int    `json:"y"`
	Width          int    `json:"width"`
	Height         int    `json:"height"`
	EncodingVector string `json:"encoding_vector"`
}
