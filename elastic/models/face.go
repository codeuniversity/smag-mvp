package models

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"math"

	"github.com/codeuniversity/smag-mvp/insta/models"
)

// FaceDoc is the type that is used to store a face in elasticsearch
type FaceDoc struct {
	PostID         int    `json:"post_id"`
	X              int    `json:"x"`
	Y              int    `json:"y"`
	Width          int    `json:"width"`
	Height         int    `json:"height"`
	EncodingVector string `json:"encoding_vector"`
}

// FaceDocFromFaceData returns a FaceDoc with an encoded `EncodingVector` given a faceData model
func FaceDocFromFaceData(faceData *models.FaceData) (*FaceDoc, error) {
	var encodingString string
	err := json.Unmarshal(faceData.Encoding.RawMessage, &encodingString)
	if err != nil {
		return nil, err
	}
	encoding := []float32{}
	err = json.Unmarshal([]byte(encodingString), &encoding)
	if err != nil {
		return nil, err
	}
	return &FaceDoc{
		PostID:         faceData.PostID,
		X:              faceData.X,
		Y:              faceData.Y,
		Width:          faceData.Width,
		Height:         faceData.Height,
		EncodingVector: EncodedVector(encoding),
	}, nil
}

// EncodedVector for the given encoding, used for searching and looking up faces in elastic search
func EncodedVector(encoding []float32) string {
	bytes := make([]byte, 0, 4*len(encoding))
	for _, a := range encoding {
		bits := math.Float32bits(a)
		b := make([]byte, 4)
		binary.BigEndian.PutUint32(b, bits)
		bytes = append(bytes, b...)
	}

	encoded := base64.StdEncoding.EncodeToString(bytes)
	return encoded
}
