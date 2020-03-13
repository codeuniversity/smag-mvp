package faceinserter

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	// necessary for gorm :pointup:
	_ "github.com/jinzhu/gorm/dialects/postgres"

	dbutils "github.com/codeuniversity/smag-mvp/db"
	"github.com/codeuniversity/smag-mvp/insta/models"
	"github.com/codeuniversity/smag-mvp/utils"
	"github.com/codeuniversity/smag-mvp/worker"
	"github.com/jinzhu/gorm/dialects/postgres"

	"github.com/segmentio/kafka-go"
)

// Inserter represents the face-encoding data inserter
type Inserter struct {
	*worker.Worker
	qReader *kafka.Reader
	db      *gorm.DB
}

// New returns an initilized inserter
func New(postgresHost, postgresPassword string, qReader *kafka.Reader) *Inserter {
	i := &Inserter{}
	i.qReader = qReader

	connectionString := fmt.Sprintf("host=%s user=postgres dbname=instascraper sslmode=disable", postgresHost)
	if postgresPassword != "" {
		connectionString += " " + "password=" + postgresPassword
	}

	db, err := gorm.Open("postgres", connectionString)
	utils.PanicIfNotNil(err)
	db.AutoMigrate(&models.FaceData{})
	i.db = db

	b := worker.Builder{}.WithName("insta_face_inserter").
		WithWorkStep(i.runStep).
		WithStopTimeout(10*time.Second).
		AddShutdownHook("qReader", qReader.Close).
		AddShutdownHook("postgres_client", db.Close)

	i.Worker = b.MustBuild()

	return i
}

// runStep the inserter
func (i *Inserter) runStep() error {
	m, err := i.qReader.FetchMessage(context.Background())
	if err != nil {
		return err
	}
	encoding := &models.FaceRecognitionResult{}
	err = json.Unmarshal(m.Value, encoding)
	if err != nil {
		return err
	}
	fmt.Println("inserting: ", encoding.PostID)
	err = i.InsertFaceEncoding(encoding)
	if err != nil {
		return err
	}

	i.qReader.CommitMessages(context.Background(), m)
	fmt.Println("commited: ", encoding.PostID)

	return nil
}

// InsertFaceEncoding inserts the encoded face data into postgres
func (i *Inserter) InsertFaceEncoding(reconResult *models.FaceRecognitionResult) error {
	p := []*models.FaceData{}
	for _, face := range reconResult.Faces {
		encodingJSON, err := json.Marshal(face.Encoding)
		if err != nil {
			return err
		}

		q := &models.FaceData{
			PostID:   reconResult.PostID,
			X:        face.X,
			Y:        face.Y,
			Width:    face.Width,
			Height:   face.Height,
			Encoding: postgres.Jsonb{RawMessage: encodingJSON},
		}
		p = append(p, q)
	}

	return i.insertEncoding(p)
}

func (i *Inserter) insertEncoding(p []*models.FaceData) error {
	fromEncoding := models.FaceData{}

	for _, face := range p {
		err := dbutils.Create(i.db, &fromEncoding, face)
		if err != nil {
			return err
		}
	}

	return nil
}
