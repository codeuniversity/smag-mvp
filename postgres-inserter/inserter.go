package inserter

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/codeuniversity/smag-mvp/utils"
	"github.com/codeuniversity/smag-mvp/worker"

	// necessary for sql :pointup:
	_ "github.com/lib/pq"

	"github.com/codeuniversity/smag-mvp/models"

	"github.com/segmentio/kafka-go"
)

// Inserter represents the scraper containing all clients it uses
type Inserter struct {
	*worker.Worker

	qReader *kafka.Reader
	qWriter *kafka.Writer
	db      *sql.DB
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

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}
	i.db = db

	b := worker.Builder{}.WithName("postgres_inserter").
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
	fmt.Println("inserting: ", info.UserName)
	err = i.InsertUserFollowInfo(info)
	if err != nil {
		return err
	}

	i.qReader.CommitMessages(context.Background(), m)
	fmt.Println("commited: ", info.UserName)

	return nil
}

// InsertUserFollowInfo inserts the user follow info into postgres
func (i *Inserter) InsertUserFollowInfo(followInfo *models.UserFollowInfo) error {
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

	return i.insertUser(p)
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func (i *Inserter) insertUser(p *models.User) error {
	_, err := i.db.Exec(`INSERT INTO users(user_name,real_name, bio, avatar_url, crawl_ts)
	VALUES($1,$2,$3,$4,$5) ON CONFLICT (user_name) DO UPDATE SET user_name = $1, real_name = $2, bio = $3, avatar_url = $4, crawl_ts = $5`,
		p.Name, p.RealName, p.Bio, p.AvatarURL, p.CrawledAt)
	if err != nil {
		return err
	}
	var userID int
	err = i.db.QueryRow("SELECT id from users where user_name = $1", p.Name).Scan(&userID)
	if err != nil {
		return err
	}

	for _, follow := range p.Follows {
		var followedID int
		err := i.db.QueryRow("SELECT id from users where user_name = $1", follow.Name).Scan(&followedID)
		if err != nil {
			err = i.db.QueryRow(`INSERT INTO users(user_name)
		VALUES($1) RETURNING id`, follow.Name).
				Scan(&followedID)
			if err != nil {
				return err
			}
			utils.HandleCreatedUser(i.qWriter, follow.Name)
		}

		_, err = i.db.Exec(`INSERT INTO follows(from_id, to_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, userID, followedID)
		if err != nil {
			return err
		}
	}
	return nil
}
