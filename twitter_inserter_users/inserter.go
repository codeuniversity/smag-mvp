package inserter

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	// necessary for gorm :pointup:
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/codeuniversity/smag-mvp/models"
	"github.com/codeuniversity/smag-mvp/utils"
	"github.com/codeuniversity/smag-mvp/worker"

	"github.com/segmentio/kafka-go"
)

// Inserter represents the scraper containing all clients it uses
type Inserter struct {
	*worker.Worker

	qReader *kafka.Reader
	qWriter *kafka.Writer

	db *gorm.DB
}

// New returns an initilized scraper
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

	db.AutoMigrate(&models.TwitterUser{})

	b := worker.Builder{}.WithName("twitter_inserter_users").
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

	err = createOrUpdate(i.db, baseUser, filter, user)
	if err != nil {
		return err
	}

	usersList := models.NewTwitterUserList(user.FollowersList, user.FollowingList)
	usersList.RemoveDuplicates()

	for _, relationUser := range *usersList {

		var queryUser models.TwitterUser

		err = i.db.Where(&models.TwitterUser{Username: relationUser.Username}).First(&queryUser).Error
		if err != nil {
			return err
		}

		if queryUser.UserIdentifier == "" {
			i.handleCreatedUser(relationUser.Username)
		}
	}
	return nil
}

func (i *Inserter) handleCreatedUser(userName string) {
	// if qWriter is nil, user discovery is disabled
	if i.qWriter != nil {
		i.qWriter.WriteMessages(context.Background(), kafka.Message{
			Value: []byte(userName),
		})
	}
}

func createOrUpdate(db *gorm.DB, out interface{}, where interface{}, update interface{}) error {
	var err error

	tx := db.Begin()

	if tx.Where(where).First(out).RecordNotFound() {
		// If the record does'nt exist it gets created
		err = tx.Create(update).Scan(out).Error
	} else {
		// Else it gets upated
		err = tx.Model(out).Update(update).Scan(out).Error
	}
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
