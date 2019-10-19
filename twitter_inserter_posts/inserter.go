package inserter

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	// necessary for sql :pointup:
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/codeuniversity/smag-mvp/models"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/codeuniversity/smag-mvp/utils"

	"github.com/segmentio/kafka-go"
)

// Inserter represents the scraper containing all clients it uses
type Inserter struct {
	qReader *kafka.Reader
	qWriter *kafka.Writer

	db *gorm.DB
	*service.Executor
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

	db.AutoMigrate(&models.TwitterPost{})

	i.Executor = service.New()
	return i
}

// Run the inserter
func (i *Inserter) Run() {
	defer func() {
		i.MarkAsStopped()
	}()

	fmt.Println("starting twitter postgres posts inserter")
	for i.IsRunning() {
		m, err := i.qReader.FetchMessage(context.Background())
		if err != nil {
			fmt.Println(err)
			break
		}
		rawPost := &models.TwitterPostRaw{}
		err = json.Unmarshal(m.Value, rawPost)
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println("inserting post:", rawPost.Link)
		var post *models.TwitterPost = models.ConvertTwitterPost(rawPost)
		i.insertPost(post)
		i.qReader.CommitMessages(context.Background(), m)
		fmt.Println("commited: ", rawPost.Link)
	}
}

// Close the inserter
func (i *Inserter) Close() {
	i.Stop()
	i.WaitUntilStopped(time.Second * 3)

	i.db.Close()
	i.qReader.Close()
	if i.qWriter != nil {
		i.qWriter.Close()
	}

	i.MarkAsClosed()
}

func (i *Inserter) insertPost(post *models.TwitterPost) {
	var err error

	fromPost := models.TwitterPost{}
	filter := &models.TwitterPost{PostIdentifier: post.PostIdentifier}

	err = createOrUpdate(i.db, &fromPost, filter, post)
	utils.PanicIfNotNil(err)

	usersList := models.NewTwitterUserList(post.Mentions, post.ReplyTo, []*models.TwitterUser{post.RetweetUser})
	usersList.RemoveDuplicates()

	for _, user := range *usersList {
		var toPost models.TwitterUser
		var d *gorm.DB
		d = i.db.Where("username = ?", user.Username).Select("ID").Find(&toPost)
		if err := d.Error; err != nil {
			if d.RecordNotFound() == true {
				d = i.db.Create(&models.TwitterUser{
					Username: user.Username,
				})
				utils.PanicIfNotNil(d.Error)

				i.handleCreatedUser(user.Username)
			} else {
				utils.PanicIfNotNil(err)
			}
		}
	}

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
		err = tx.Create(update).Scan(out).Error
	} else {
		err = tx.Model(out).Update(update).Scan(out).Error
	}
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}
