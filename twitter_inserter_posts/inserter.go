package inserter

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	// necessary for gorm :pointup:
	_ "github.com/jinzhu/gorm/dialects/postgres"

	dbUtils "github.com/codeuniversity/smag-mvp/db"
	"github.com/codeuniversity/smag-mvp/models"
	"github.com/codeuniversity/smag-mvp/utils"
	"github.com/codeuniversity/smag-mvp/worker"

	"github.com/segmentio/kafka-go"
)

// Inserter represents the inserter containing all clients it uses
type Inserter struct {
	*worker.Worker

	qReader *kafka.Reader

	db *gorm.DB
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
	i.db = db

	db.AutoMigrate(&models.TwitterPost{})

	b := worker.Builder{}.WithName("twitter_inserter_posts").
		WithWorkStep(i.runStep).
		WithStopTimeout(10*time.Second).
		AddShutdownHook("qReader", qReader.Close).
		AddShutdownHook("postgres_client", db.Close)

	i.Worker = b.MustBuild()

	return i
}

// Run the inserter
func (i *Inserter) runStep() error {
	m, err := i.qReader.FetchMessage(context.Background())
	if err != nil {
		return err
	}

	rawPost := &models.TwitterPostRaw{}

	err = json.Unmarshal(m.Value, rawPost)
	if err != nil {
		return err
	}

	post := models.ConvertTwitterPost(rawPost)
	fmt.Println("inserting post:", post.Link)

	err = i.insertPost(post)
	if err != nil {
		return err
	}
	return i.qReader.CommitMessages(context.Background(), m)
}

func (i *Inserter) insertPost(post *models.TwitterPost) error {
	fromPost := &models.TwitterPost{}
	filter := &models.TwitterPost{PostIdentifier: post.PostIdentifier}

	err := dbUtils.CreateOrUpdate(i.db, fromPost, filter, post)
	if err != nil {
		return err
	}

	return nil
}
