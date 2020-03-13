package inserter

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	// necessary for gorm :pointup:
	_ "github.com/jinzhu/gorm/dialects/postgres"

	dbUtils "github.com/codeuniversity/smag-mvp/db"
	"github.com/codeuniversity/smag-mvp/insta/models"
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
	db.AutoMigrate(&models.User{}, &models.Follow{})
	i.db = db

	b := worker.Builder{}.WithName("insta_postgres_inserter").
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
	info := &models.UserFollowInfo{}
	err = json.Unmarshal(m.Value, info)
	if err != nil {
		return err
	}
	log.Println("inserting: ", info.UserName)
	err = i.InsertUserFollowInfo(info)
	if err != nil {
		return err
	}

	i.qReader.CommitMessages(context.Background(), m)
	log.Println("commited: ", info.UserName)

	return nil
}

// InsertUserFollowInfo inserts the user follow info into postgres
func (i *Inserter) InsertUserFollowInfo(followInfo *models.UserFollowInfo) error {
	p := &models.User{
		UserName:  followInfo.UserName,
		RealName:  followInfo.RealName,
		AvatarURL: followInfo.AvatarURL,
		Bio:       followInfo.Bio,
		CrawledAt: followInfo.CrawlTs,
	}

	for _, following := range followInfo.Followings {
		p.Follows = append(p.Follows, &models.User{
			UserName: following,
		})
	}

	return i.insertUser(p)
}

func (i *Inserter) insertUser(p *models.User) error {
	fromUser := models.User{}
	filter := &models.User{UserName: p.UserName}

	err := dbUtils.CreateOrUpdate(i.db, &fromUser, filter, p)
	if err != nil {
		return err
	}

	for _, follow := range p.Follows {
		var toUser models.User
		var d *gorm.DB
		d = i.db.Where("user_name = ?", follow.UserName).Select("ID").Find(&toUser)
		if err := d.Error; err != nil {
			if d.RecordNotFound() == true {
				d = i.db.Create(&models.User{
					UserName: follow.UserName,
				}).Scan(&toUser)
				if d.Error != nil {
					return d.Error
				}

			} else {
				if err != nil {
					return err
				}
			}
		}

		err = i.db.Set("gorm:insert_option", "ON CONFLICT DO NOTHING").Create(&models.Follow{
			From: fromUser.ID,
			To:   toUser.ID,
		}).Error
		if err != nil {
			return err
		}
	}
	return nil
}
