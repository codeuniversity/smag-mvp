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

	db.AutoMigrate(&models.TwitterUser{})

	i.Executor = service.New()
	return i
}

// Run the inserter
func (i *Inserter) Run() {
	defer i.MarkAsStopped()

	fmt.Println("starting inserter")
	for i.IsRunning() {
		m, err := i.qReader.FetchMessage(context.Background())
		if err != nil {
			fmt.Println(err)
			break
		}
		info := &models.TwitterUserRaw{}
		err = json.Unmarshal(m.Value, info)
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println("inserting user: ", info.Username)
		i.insertUser(models.ConvertTwitterUser(info))
		i.qReader.CommitMessages(context.Background(), m)
		fmt.Println("commited: ", info.Username)
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

func (i *Inserter) insertUser(user *models.TwitterUser) {
	var err error

	baseUser := models.TwitterUser{}
	filter := &models.TwitterUser{Username: user.Username}

	err = createOrUpdate(i.db, &baseUser, filter, user)
	utils.PanicIfNotNil(err)

	usersList := models.NewTwitterUserList(user.FollowersList, user.FollowingList)
	usersList.RemoveDuplicates()

	for _, relationUser := range *usersList {
		fmt.Printf("\nHandle user %v\n%v\n", relationUser.Username, relationUser)

		var queryUser models.TwitterUser

		err := i.db.Where("username = ?", relationUser.Username).Find(&queryUser).Error
		utils.PanicIfNotNil(err)

		fmt.Println("Query resulted in:\n", queryUser)
		if queryUser.Username == "" {
			fmt.Printf("Didn't find record for %v\n", relationUser.Username)
			err = i.db.Create(&models.TwitterUser{
				Username: relationUser.Username,
			}).Error
			utils.PanicIfNotNil(err)

			i.handleCreatedUser(relationUser.Username)
		}

		fmt.Println("Done handling.")
	}
}

func (i *Inserter) handleCreatedUser(userName string) {
	// if qWriter is nil, user discovery is disabled
	if i.qWriter != nil {
		fmt.Printf("Send %v to kafka/%v\n", userName, i.qWriter.Stats().Topic)
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
		fmt.Println("Insert user ", update, " into postgres")
		err = tx.Create(update).Scan(out).Error
	} else {
		// Else it gets upated
		fmt.Println("Update user ", update, " in postgres")
		err = tx.Model(out).Update(update).Scan(out).Error
	}
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
