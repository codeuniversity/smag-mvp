package detection

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	http_header_generator "github.com/codeuniversity/smag-mvp/http_header-generator"
	"github.com/codeuniversity/smag-mvp/insta/models"
	"github.com/codeuniversity/smag-mvp/utils"
	"github.com/codeuniversity/smag-mvp/worker"

	"github.com/minio/minio-go"
	"github.com/segmentio/kafka-go"
	"gocv.io/x/gocv"
)

// Detector represents the detector containing all clients it uses
type Detector struct {
	*worker.Worker

	nameQReader *kafka.Reader
	infoQWriter *kafka.Writer
	errQWriter  *kafka.Writer
	*http_header_generator.HTTPHeaderGenerator

	minioClient *minio.Client
	bucketName  string
	region      string

	classifier gocv.CascadeClassifier
}

// Config holds all s3 related configs
type Config struct {
	S3BucketName      string
	S3Region          string
	S3Endpoint        string
	S3AccessKeyID     string
	S3SecretAccessKey string
	S3UseSSL          bool
}

// New returns an initilized scraper
func New(nameQReader *kafka.Reader, infoQWriter *kafka.Writer, config Config) *Detector {
	d := &Detector{}
	d.nameQReader = nameQReader
	d.infoQWriter = infoQWriter
	d.HTTPHeaderGenerator = http_header_generator.New()

	d.bucketName = config.S3BucketName
	d.region = config.S3Region

	minioClient, err := minio.New(config.S3Endpoint, config.S3AccessKeyID, config.S3SecretAccessKey, config.S3UseSSL)
	utils.MustBeNil(err)
	d.minioClient = minioClient

	xmlFileFrontal := "haarcascade_frontalface_alt.xml"

	// load classifier to recognize faces
	d.classifier = gocv.NewCascadeClassifier()

	if !d.classifier.Load(xmlFileFrontal) {
		panic("error reading cascade file")
	}

	d.Worker = worker.Builder{}.WithName("insta_posts_face-detector").
		WithWorkStep(d.runStep).
		AddShutdownHook("nameQReader", nameQReader.Close).
		AddShutdownHook("infoQWriter", infoQWriter.Close).
		AddShutdownHook("classifier", d.classifier.Close).
		MustBuild()
	return d
}

func (d *Detector) runStep() error {
	fmt.Println("fetching")
	m, err := d.nameQReader.FetchMessage(context.Background())
	if err != nil {
		return err
	}

	job := models.PostDownloadJob{}
	err = json.Unmarshal(m.Value, &job)
	if err != nil {
		return err
	}

	if job.PictureURL != "" {
		err = d.fetchPost(job.PictureURL, job.PostID)
		if err != nil {
			return err
		}
	}

	return d.nameQReader.CommitMessages(context.Background(), m)
}

func (d *Detector) fetchPost(internalImgURL string, postID int) error {
	localImagePath := internalImgURL
	err := d.minioClient.FGetObject(d.bucketName, internalImgURL, localImagePath, minio.GetObjectOptions{})
	if err != nil {
		return err
	}
	return d.analyzeForFaces(localImagePath, postID)
}

func (d *Detector) analyzeForFaces(localImagePath string, postID int) error {

	fmt.Printf("Reading image: %s \n", localImagePath)
	fmt.Printf("Image is form post: %v \n", postID)

	img := gocv.IMRead(localImagePath, gocv.IMReadColor)
	if img.Empty() {
		fmt.Println("error reading image - image is empty")
		return nil
	}

	picture, err := img.ToImage()
	if err != nil {
		panic(err)
	}

	defer func() {
		img.Close()
		err := os.Remove(localImagePath)
		if err != nil {
			panic(err)
		}
	}()
	// detect frontals
	rects := d.classifier.DetectMultiScale(img)

	if len(rects) == 0 {
		return nil
	}

	faceReconJob := &models.FaceReconJob{}
	faceReconJob.PostID = postID
	faceReconJob.InternalImageURL = localImagePath

	if len(rects) >= 1 {
		faceReconJob.X = 0
		faceReconJob.Y = 0
		faceReconJob.Width = picture.Bounds().Dx()
		faceReconJob.Height = picture.Bounds().Dy()
	}

	return d.writeFilteredResult(faceReconJob)
}

func (d *Detector) writeFilteredResult(job *models.FaceReconJob) error {
	payload, err := json.Marshal(job)
	if err != nil {
		return err
	}

	return d.infoQWriter.WriteMessages(context.Background(), kafka.Message{Value: payload})
}
