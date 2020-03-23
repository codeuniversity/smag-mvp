package recognition

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/codeuniversity/smag-mvp/faces/proto"
	"github.com/codeuniversity/smag-mvp/imgproxy"
	"github.com/codeuniversity/smag-mvp/insta/models"
	"github.com/codeuniversity/smag-mvp/worker"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
)

// Worker reads from the face recognition kafka queue,
// calls the face_recognition grpc service and writes the found faces into another queue
type Worker struct {
	jobQReader            *kafka.Reader
	resultQWriter         *kafka.Writer
	faceRecognizerAddress string
	urlBuilder            *imgproxy.URLBuilder
	bucketName            string
	*worker.Worker
}

// New returns an intialized worker
func New(jobQReader *kafka.Reader, resultQWriter *kafka.Writer, faceRecognizerAddress string, pictureBucketName string, imgProxyAddress, imgProxyKey, imgProxySalt string) *Worker {
	urlBuilder, err := imgproxy.New(imgProxyAddress, imgProxyKey, imgProxySalt)
	if err != nil {
		panic(err)
	}
	w := &Worker{
		jobQReader:            jobQReader,
		resultQWriter:         resultQWriter,
		urlBuilder:            urlBuilder,
		bucketName:            pictureBucketName,
		faceRecognizerAddress: faceRecognizerAddress,
	}

	w.Worker = worker.Builder{}.WithName("face_recognition_worker").
		WithWorkStep(w.step).
		AddShutdownHook("jobQReader", jobQReader.Close).
		AddShutdownHook("resultQWriter", resultQWriter.Close).
		WithStopTimeout(20 * time.Second).
		MustBuild()

	return w
}

func (w *Worker) step() error {
	messages, err := w.readMessageBlock(1*time.Second, 8)
	if err != nil {
		return err
	}

	if len(messages) == 0 {
		return nil
	}

	resultChan := make(chan recognitionResult)
	for _, message := range messages {
		go func(resultChan chan recognitionResult, jobMessage kafka.Message) {
			job := &models.FaceReconJob{}
			err = json.Unmarshal(jobMessage.Value, job)
			if err != nil {
				resultChan <- recognitionResult{Error: err}
				return
			}

			if strings.Trim(job.InternalImageURL, " ") == "" {
				log.Println("Empty image URL for post", job.PostID)
				resultChan <- recognitionResult{JobMessage: jobMessage}
				return
			}

			con, err := grpc.Dial(w.faceRecognizerAddress, grpc.WithInsecure())
			if err != nil {
				resultChan <- recognitionResult{Error: err}
				return
			}
			defer con.Close()

			client := proto.NewFaceRecognizerClient(con)

			url := w.urlBuilder.GetCropURL(0, 0, 50000, 50000, w.urlBuilder.GetS3Url(w.bucketName, job.InternalImageURL))
			response, err := client.RecognizeFaces(context.Background(), &proto.RecognizeRequest{
				Url: url,
			})
			if err != nil {
				resultChan <- recognitionResult{Error: err}
				return
			}
			faces := response.Faces
			if len(faces) == 0 {
				resultChan <- recognitionResult{JobMessage: jobMessage}
				return
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
					err := fmt.Errorf("face encoding has wrong len, expected 128 but was %d", len(face.Encoding))
					resultChan <- recognitionResult{Error: err}
					return
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
				resultChan <- recognitionResult{Error: err}
				return
			}
			resultChan <- recognitionResult{ResultMesage: &kafka.Message{Value: b}, JobMessage: jobMessage}
		}(resultChan, message)
	}

	collectedErrors := []error{}
	collectedWrites := []kafka.Message{}
	messagesToCommit := []kafka.Message{}
	for i := 0; i < len(messages); i++ {
		r := <-resultChan

		if r.Error != nil {
			collectedErrors = append(collectedErrors, r.Error)
			continue
		}

		messagesToCommit = append(messagesToCommit, r.JobMessage)
		if r.ResultMesage != nil {
			collectedWrites = append(collectedWrites, *r.ResultMesage)
		}
	}

	if len(collectedErrors) > 0 {
		mergedError := errors.New("not all recognition requests were successful")
		for _, err := range collectedErrors {
			mergedError = fmt.Errorf("%v, %w", mergedError, err)
		}

		return mergedError
	}

	err = w.resultQWriter.WriteMessages(context.Background(), collectedWrites...)
	if err != nil {
		return err
	}

	return w.jobQReader.CommitMessages(context.Background(), messagesToCommit...)
}

type recognitionResult struct {
	JobMessage   kafka.Message
	ResultMesage *kafka.Message
	Error        error
}

func (w *Worker) readMessageBlock(timeout time.Duration, maxChunkSize int) (messages []kafka.Message, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for i := 0; i < maxChunkSize; i++ {
		m, err := w.jobQReader.FetchMessage(ctx)
		if err != nil {
			if err == context.DeadlineExceeded {
				return messages, nil
			}

			return nil, err
		}

		messages = append(messages, m)
	}

	return messages, nil
}
