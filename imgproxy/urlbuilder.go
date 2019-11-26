package imgproxy

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
)

// URLBuilder simplifies constructing urls for the imgproxy
type URLBuilder struct {
	keyBin       []byte
	saltBin      []byte
	proxyAddress string
}

// New returns a URLBuilder or an error if the salt or key is not hex-encoded
func New(proxyAddress, key, salt string) (*URLBuilder, error) {
	var keyBin, saltBin []byte
	var err error

	if keyBin, err = hex.DecodeString(key); err != nil {
		return nil, errors.New("Key expected to be hex-encoded string")
	}

	if saltBin, err = hex.DecodeString(salt); err != nil {
		return nil, errors.New("Salt expected to be hex-encoded string")
	}

	return &URLBuilder{
		keyBin:       keyBin,
		saltBin:      saltBin,
		proxyAddress: proxyAddress,
	}, nil
}

// GetCropURL returns a url that instructs the imgproxy to crop out the given cordinates of the image
// with a sourceURL in the form "s3://<bucket_name>/<path_to_image>" the image proxy will download the image from s3
func (b *URLBuilder) GetCropURL(x, y, width, height int, sourceURL string) string {
	encodedURL := base64.RawURLEncoding.EncodeToString([]byte(sourceURL))
	gravity := fmt.Sprintf("nowe:%d:%d", x, y)
	extension := "jpg"
	path := fmt.Sprintf("/crop:%d:%d:%s/%s.%s", width, height, gravity, encodedURL, extension)

	mac := hmac.New(sha256.New, b.keyBin)
	mac.Write(b.saltBin)
	mac.Write([]byte(path))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))[:32]

	return fmt.Sprintf("http://%s/%s%s", b.proxyAddress, signature, path)
}

// GetS3Url returns a url in the form "s3://<bucket_name>/<path_to_image>"
func (b *URLBuilder) GetS3Url(bucketName, path string) string {
	return fmt.Sprintf("s3://%s/%s", bucketName, path)
}
