package models

// Config holds all the configurable variables for the Downloader
type Config struct {
	S3BucketName      string
	S3Region          string
	S3Endpoint        string
	S3AccessKeyID     string
	S3SecretAccessKey string
	S3UseSSL          bool

	PostgresHost     string
	PostgresPassword string
}
