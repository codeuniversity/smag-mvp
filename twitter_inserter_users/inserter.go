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

// Inserter represents the scraper containing all clients it uses
type Inserter struct {
	*worker.Worker

	qReader *kafka.Reader

	db *gorm.DB
}

// New returns an initilized scraper
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

	db.AutoMigrate(&models.TwitterUser{})

	b := worker.Builder{}.WithName("twitter_inserter_users").
		WithWorkStep(i.runStep).
		WithStopTimeout(10*time.Second).
		AddShutdownHook("qReader", qReader.Close).
		AddShutdownHook("postgres_client", db.Close)

	i.Worker = b.MustBuild()

	return i
}

func (i *Inserter) runStep() error {
	m, err := i.qReader.FetchMessage(context.Background())
	if err != nil {
		return err
	}

	rawUser := &models.TwitterUserRaw{}

	err = json.Unmarshal(m.Value, rawUser)
	if err != nil {
		return err
	}

	user := models.ConvertTwitterUser(rawUser)
	fmt.Println("inserting user: ", user.Username)

	err = i.insertUser(user)
	if err != nil {
		return err
	}

	return i.qReader.CommitMessages(context.Background(), m)
}

func (i *Inserter) insertUser(user *models.TwitterUser) error {
	var err error

	baseUser := &models.TwitterUser{}
	filter := &models.TwitterUser{Username: user.Username}

	err = dbUtils.CreateOrUpdate(i.db, baseUser, filter, user)
	if err != nil {
		return err
	}

	return nil
}
