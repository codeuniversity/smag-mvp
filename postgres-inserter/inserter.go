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
func New(kafkaAddress, postgresHost, postgresPassword string) *Inserter {
	i := &Inserter{}
	i.qReader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{kafkaAddress},
		GroupID:        "user_postgres_inserter",
		Topic:          "user_follow_infos",
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		CommitInterval: time.Second,
	})
	i.qWriter = kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{kafkaAddress},
		Topic:    "user_names",
		Balancer: &kafka.LeastBytes{},
		Async:    true,
	})
	connectionString := fmt.Sprintf("host=%s user=postgres dbname=instascraper sslmode=disable", postgresHost)
	if postgresPassword != "" {
		connectionString += " " + "password=" + postgresPassword
	}

	db, err := gorm.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}
	i.db = db

	db.AutoMigrate(&models.User{})

	i.Executor = service.New()
	return i
}

// Run the inserter
func (i *Inserter) Run() {
	defer func() {
		i.MarkAsStopped()
	}()

	fmt.Println("starting inserter")
	for i.IsRunning() {
		m, err := i.qReader.FetchMessage(context.Background())
		if err != nil {
			fmt.Println(err)
			break
		}
		info := &models.UserFollowInfo{}
		err = json.Unmarshal(m.Value, info)
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println("inserting: ", info.UserName)
		i.InsertUserFollowInfo(info)
		i.qReader.CommitMessages(context.Background(), m)
		fmt.Println("commited: ", info.UserName)
	}
}

// Close the inserter
func (i *Inserter) Close() {
	i.Stop()
	i.WaitUntilStopped(time.Second * 3)

	i.db.Close()
	i.qReader.Close()
	i.qWriter.Close()

	i.MarkAsClosed()
}

// InsertUserFollowInfo inserts the user follow info into dgraph, while writting userNames that don't exist in the graph yet
// into the specified kafka topic
func (i *Inserter) InsertUserFollowInfo(followInfo *models.UserFollowInfo) {
	p := &models.User{
		Name:      followInfo.UserName,
		RealName:  followInfo.RealName,
		AvatarURL: followInfo.AvatarURL,
		Bio:       followInfo.Bio,
		CrawledAt: followInfo.CrawlTs,
	}

	for _, following := range followInfo.Followings {
		p.Follows = append(p.Follows, &models.User{
			Name: following,
		})
	}

	i.insertUser(p)
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func (i *Inserter) insertUser(p *models.User) {
	/*_, err := i.db.Exec(`INSERT INTO users(user_name,real_name, bio, avatar_url, crawl_ts)
	VALUES($1,$2,$3,$4,$5) ON CONFLICT (user_name) DO UPDATE SET user_name = $1, real_name = $2, bio = $3, avatar_url = $4, crawl_ts = $5`,
		p.Name, p.RealName, p.Bio, p.AvatarURL, p.CrawledAt)
	*/
	var err error

	fromUser := models.User{}
	filter := &models.User{Name: p.Name}

	err = createOrUpdate(i.db, &fromUser, filter, p)
	handleErr(err)

	for _, follow := range p.Follows {
		//err := i.db.QueryRow("SELECT id from users where user_name = $1", follow.Name).Scan(&followedID)
		var toUser models.User
		var d *gorm.DB
		d = i.db.Where("user_name = ?", follow.Name).Select("ID").Find(&toUser)
		if err := d.Error; err != nil {
			// err = i.db.QueryRow(`INSERT INTO users(user_name) VALUES($1) RETURNING id`, follow.Name).Scan(&followedID)
			if d.RecordNotFound() == true {
				d = i.db.Create(&models.User{
					Name: follow.Name,
				}).Scan(&toUser)
				handleErr(d.Error)

				i.handleCreatedUser(follow.Name)
			} else {
				handleErr(err)
			}
		}

		//_, err = i.db.Exec(`INSERT INTO follows(from_id, to_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, userID, followedID)
		err = i.db.Set("gorm:insert_option", "ON CONFLICT DO NOTHING").Create(&models.Follow{
			From: fromUser.ID,
			To:   toUser.ID,
		}).Error
		handleErr(err)
	}

}

func (i *Inserter) handleCreatedUser(userName string) {
	i.qWriter.WriteMessages(context.Background(), kafka.Message{
		Value: []byte(userName),
	})
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
