package config

import "github.com/codeuniversity/smag-mvp/utils"

// S3Config holds all the configurable variables for S3
type S3Config struct {
	S3BucketName      string
	S3Region          string
	S3Endpoint        string
	S3AccessKeyID     string
	S3SecretAccessKey string
	S3UseSSL          bool
}

//GetS3Config returns a inizialized S3 Config
func GetS3Config() *S3Config {
	return &S3Config{
		S3BucketName:      utils.GetStringFromEnvWithDefault("S3_BUCKET_NAME", "insta_pics"),
		S3Region:          utils.GetStringFromEnvWithDefault("S3_REGION", "eu-west-1"),
		S3Endpoint:        utils.GetStringFromEnvWithDefault("S3_ENDOINT", "127.0.0.1:9000"),
		S3AccessKeyID:     utils.MustGetStringFromEnv("S3_ACCESS_KEY_ID"),
		S3SecretAccessKey: utils.MustGetStringFromEnv("S3_SECRET_ACCESS_KEY"),
		S3UseSSL:          utils.GetBoolFromEnvWithDefault("S3_USE_SSL", true),
	}
}
