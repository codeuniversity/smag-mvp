package models

// FaceRecognitionResult is the result of the face recognizer
type FaceRecognitionResult struct {
	PostID int     `json:"post_id"`
	Faces  []*Face `json:"faces"`
}

// Face contains the position of a face in a post and its encoding
type Face struct {
	X        int          `json:"x"`
	Y        int          `json:"y"`
	Width    int          `json:"width"`
	Height   int          `json:"height"`
	Encoding [128]float32 `json:"encoding"`
}
