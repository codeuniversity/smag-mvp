package recognition

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/codeuniversity/smag-mvp/faces/proto"
	"github.com/codeuniversity/smag-mvp/imgproxy"
	"github.com/codeuniversity/smag-mvp/models"
	"github.com/codeuniversity/smag-mvp/worker"
	kgo "github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
)

// Worker reads from the face recognition kafka queue,
// calls the face_recognition grpc service and writes the found faces into another queue
type Worker struct {
	jobQReader       *kgo.Reader
	resultQWriter    *kgo.Writer
	recognizerClient proto.FaceRecognizerClient
	urlBuilder       *imgproxy.URLBuilder
	bucketName       string
	*worker.Worker
}

// New returns an intialized worker
func New(jobQReader *kgo.Reader, resultQWriter *kgo.Writer, faceRecognizerAddress string, pictureBucketName string, imgProxyAddress, imgProxyKey, imgProxySalt string) *Worker {
	con, err := grpc.Dial(faceRecognizerAddress, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client := proto.NewFaceRecognizerClient(con)
	urlBuilder, err := imgproxy.New(imgProxyAddress, imgProxyKey, imgProxySalt)
	if err != nil {
		panic(err)
	}
	w := &Worker{
		jobQReader:       jobQReader,
		resultQWriter:    resultQWriter,
		recognizerClient: client,
		urlBuilder:       urlBuilder,
		bucketName:       pictureBucketName,
	}

	w.Worker = worker.Builder{}.WithName("face_recognition_worker").
		WithWorkStep(w.step).
		AddShutdownHook("jobQReader", jobQReader.Close).
		AddShutdownHook("resultQWriter", resultQWriter.Close).
		AddShutdownHook("face_recognizer_connection", con.Close).
		MustBuild()

	return w
}

func (w *Worker) step() error {
	m, err := w.jobQReader.FetchMessage(context.Background())
	if err != nil {
		return err
	}

	job := &models.FaceReconJob{}
	err = json.Unmarshal(m.Value, job)
	if err != nil {
		return err
	}
	url := w.urlBuilder.GetCropURL(job.X, job.Y, job.Width, job.Height, w.urlBuilder.GetS3Url(w.bucketName, job.InternalImageURL))
	response, err := w.recognizerClient.RecognizeFaces(context.Background(), &proto.RecognizeRequest{
		Url: url,
	})
	if err != nil {
		return err
	}
	faces := response.Faces
	if len(faces) == 0 {
		return w.jobQReader.CommitMessages(context.Background(), m)
	}
	result := &models.FaceRecognitionResult{PostID: job.PostID}
	for _, face := range faces {
		x := int(face.X)
		y := int(face.Y)
		width := int(face.Width)
		height := int(face.Height)
		url := w.urlBuilder.GetCropURL(x, y, width, height, w.urlBuilder.GetS3Url(w.bucketName, job.InternalImageURL))
		fmt.Println(url)

		if len(face.Encoding) != 128 {
			log.Fatal("face encoding has wrong len, expected 128 but was", len(face.Encoding))
		}

		encoding := [128]float32{}

		for i, v := range face.Encoding {
			encoding[i] = v
		}

		resultFace := &models.Face{
			X:        x,
			Y:        y,
			Width:    width,
			Height:   height,
			Encoding: encoding,
		}

		result.Faces = append(result.Faces, resultFace)
	}

	b, err := json.Marshal(result)
	if err != nil {
		return err
	}

	err = w.resultQWriter.WriteMessages(context.Background(), kgo.Message{Value: b})
	if err != nil {
		return err
	}

	return w.jobQReader.CommitMessages(context.Background(), m)
}
