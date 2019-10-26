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
	qWriter *kafka.Writer

	db *gorm.DB
}

// New returns an initilized inserter
func New(postgresHost, postgresPassword string, qReader *kafka.Reader, qWriter *kafka.Writer) *Inserter {
	i := &Inserter{}
	i.qReader = qReader
	i.qWriter = qWriter

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

	if qWriter != nil {
		b = b.AddShutdownHook("qWriter", qWriter.Close)
	}

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

	newUserLists := [][]*models.TwitterUser{post.Mentions, post.ReplyTo}
	if post.RetweetUserID != "" {
		newUserLists = append(newUserLists, []*models.TwitterUser{
			&models.TwitterUser{
				Username: post.RetweetUsername,
			},
		})
	}
	usersList := models.NewTwitterUserList(newUserLists...)
	usersList.RemoveDuplicates()

	for _, relationUser := range *usersList {

		queryUser := models.TwitterUser{}

		err = i.db.Where("username = ?", relationUser.Username).First(&queryUser).Error
		if err != nil {
			return err
		}
	}

	return nil
}
