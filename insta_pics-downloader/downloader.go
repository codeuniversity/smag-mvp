package downloader

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	// necessary for sql :pointup:
	_ "github.com/lib/pq"

	"github.com/codeuniversity/smag-mvp/models"
	"github.com/codeuniversity/smag-mvp/utils"
	"github.com/codeuniversity/smag-mvp/worker"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v6"
	"github.com/segmentio/kafka-go"
)

// Downloader reads download jobs from kafka, downloads pictures from posts, stores them in S3
//  and writes the S3 path to posts
type Downloader struct {
	*worker.Worker

	qReader *kafka.Reader

	minioClient *minio.Client
	bucketName  string
	region      string
	db          *sql.DB
}

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

// New returns an initilized scraper
func New(qReader *kafka.Reader, config Config) *Downloader {
	i := &Downloader{}

	i.bucketName = config.S3BucketName
	i.region = config.S3Region

	minioClient, err := minio.New(config.S3Endpoint, config.S3AccessKeyID, config.S3SecretAccessKey, config.S3UseSSL)
	utils.MustBeNil(err)

	i.minioClient = minioClient
	err = i.ensureBucketExists()
	utils.MustBeNil(err)

	i.qReader = qReader

	connectionString := fmt.Sprintf("host=%s user=postgres dbname=instascraper sslmode=disable", config.PostgresHost)
	if config.PostgresPassword != "" {
		connectionString += " " + "password=" + config.PostgresPassword
	}

	db, err := sql.Open("postgres", connectionString)
	utils.MustBeNil(err)
	i.db = db

	i.Worker = worker.Builder{}.WithName("insta_pics_downloader").
		WithWorkStep(i.runStep).
		WithStopTimeout(10*time.Second).
		AddShutdownHook("postgres_connection", db.Close).
		AddShutdownHook("qReader", i.qReader.Close).
		MustBuild()

	return i
}

func (d *Downloader) runStep() error {
	m, err := d.qReader.FetchMessage(context.Background())
	if err != nil {
		return err
	}

	job := models.PostDownloadJob{}
	err = json.Unmarshal(m.Value, &job)
	if err != nil {
		return err
	}

	path, err := d.downloadImgToS3(job)
	if err != nil {
		return err
	}

	err = d.updatePost(job.PostID, path)
	if err != nil {
		return err
	}

	return d.qReader.CommitMessages(context.Background(), m)
}

func (d *Downloader) downloadImgToS3(job models.PostDownloadJob) (path string, err error) {
	response, err := http.Get(job.PictureURL)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	path = randSeq()
	n, err := d.minioClient.PutObject(d.bucketName, path, response.Body, response.ContentLength, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		return "", err
	}
	log.Println("downloaded", job.PictureURL, "size:", n, "to", path)
	return path, nil

}

func (d *Downloader) updatePost(postID int, internalPath string) error {
	_, err := d.db.Exec("UPDATE posts SET internal_picture_url = $1 where id = $2", internalPath, postID)
	return err
}

func (d *Downloader) ensureBucketExists() error {
	exists, err := d.minioClient.BucketExists(d.bucketName)
	if err != nil {
		return err
	}

	if exists {
		log.Println("Bucket", d.bucketName, "already exists")
		return nil
	}

	err = d.minioClient.MakeBucket(d.bucketName, d.region)
	if err != nil {
		return fmt.Errorf("couldn't create bucket %s: %w", d.bucketName, err)
	}
	log.Printf("successfully created bucket %s\n", d.bucketName)
	return nil
}

func randSeq() string {
	id, err := uuid.NewRandom()
	utils.MustBeNil(err)

	return id.String()
}
