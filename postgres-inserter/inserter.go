package inserter

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	// necessary for sql :pointup:
	_ "github.com/lib/pq"

	"github.com/codeuniversity/smag-mvp/models"
	"github.com/codeuniversity/smag-mvp/service"

	"github.com/segmentio/kafka-go"
)

// Inserter represents the scraper containing all clients it uses
type Inserter struct {
	qReader *kafka.Reader
	qWriter *kafka.Writer

	db *sql.DB
	*service.Executor
}

// New returns an initilized scraper
func New(kafkaAddress, postgresHost string) *Inserter {
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
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s user=postgres dbname=instascraper sslmode=disable", postgresHost))
	if err != nil {
		panic(err)
	}
	i.db = db

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
	_, err := i.db.Exec(`INSERT INTO users(user_name,real_name, bio, avatar_url, crawl_ts)
	VALUES($1,$2,$3,$4,$5) ON CONFLICT (user_name) DO UPDATE SET user_name = $1, real_name = $2, bio = $3, avatar_url = $4, crawl_ts = $5`,
		p.Name, p.RealName, p.Bio, p.AvatarURL, p.CrawledAt)
	handleErr(err)
	var userID int
	i.db.QueryRow("SELECT id from users where user_name = $1", p.Name).Scan(&userID)

	for _, follow := range p.Follows {
		var followedID int
		err := i.db.QueryRow("SELECT id from users where user_name = $1", follow.Name).Scan(&followedID)
		if err != nil {
			err = i.db.QueryRow(`INSERT INTO users(user_name)
		VALUES($1) RETURNING id`, follow.Name).
				Scan(&followedID)
			handleErr(err)
			i.handleCreatedUser(follow.Name)
		}

		_, err = i.db.Exec(`INSERT INTO follows(from_id, to_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, userID, followedID)
		handleErr(err)
	}
}

func (i *Inserter) handleCreatedUser(userName string) {
	i.qWriter.WriteMessages(context.Background(), kafka.Message{
		Value: []byte(userName),
	})
}
