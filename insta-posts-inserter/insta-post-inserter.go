package insta_posts_inserter

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/codeuniversity/smag-mvp/models"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/segmentio/kafka-go"
	"time"
)

type InstaPostInserter struct {
	postQReader *kafka.Reader
	errQWriter  *kafka.Writer
	*service.Executor
	db           *sql.DB
	kafkaAddress string
}

func New(kafkaAddress string, postgresHost string) *InstaPostInserter {
	p := &InstaPostInserter{}
	p.postQReader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{kafkaAddress},
		GroupID:        "insta_post_inserter_group1",
		Topic:          "user_post",
		CommitInterval: time.Minute * 40,
	})
	p.errQWriter = kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{kafkaAddress},
		Topic:    "post_post_inserter_errors",
		Balancer: &kafka.LeastBytes{},
		Async:    false,
	})
	p.Executor = service.New()
	p.kafkaAddress = kafkaAddress

	db, err := sql.Open("postgres", fmt.Sprintf("host=%s user=postgres dbname=instascraper sslmode=disable", postgresHost))
	if err != nil {
		panic(err)
	}
	p.db = db
	return p
}

func (i *InstaPostInserter) findOrCreateUser(username string) (userID int, err error) {
	err = i.db.QueryRow("Select id from users where user_name = $1", username).Scan(&userID)

	if err != nil {
		if err != sql.ErrNoRows {
			return 0, err
		}

		result, err := i.db.Exec(`INSERT INTO users(user_name) VALUES($1) RETURNING id`, username)
		if err != nil {
			return 0, err
		}

		insertedUserID, err := result.LastInsertId()

		if err != nil {
			return 0, err
		}

		userID = int(insertedUserID)
	}

	return userID, nil
}

func (i *InstaPostInserter) Run() {
	defer func() {
		i.MarkAsStopped()
	}()

	fmt.Println("starting inserter")
	for i.IsRunning() {

		message, err := i.postQReader.FetchMessage(context.Background())
		if err != nil {
			fmt.Println(err)
			break
		}

		var post models.InstagramPost
		err = json.Unmarshal(message.Value, &post)
		if err != nil {
			panic(err)
		}

		err = i.insertPost(post)

		if err != nil {
			panic(fmt.Errorf("comments inserter failed %s ", err))
		}
		i.postQReader.CommitMessages(context.Background(), message)
	}
}

func (i *InstaPostInserter) insertPost(post models.InstagramPost) error {

	_, err := i.findOrCreateUser(post.UserId)

	if err != nil {
		return err
	}

	_, err = i.db.Exec(`INSERT INTO posts(user_id, post_id, short_code, picture_url) VALUES($1,$2,$3,$4) ON CONFLICT(post_id) DO UPDATE SET short_code=$2, picture_url=$4`, post.UserId, post.PostId, post.ShortCode, post.PictureUrl)

	if err != nil {
		return err
	}

	return nil
}

func (i *InstaPostInserter) Close() {
	i.Stop()
	i.WaitUntilStopped(time.Second * 3)

	i.errQWriter.Close()
	i.postQReader.Close()
	i.MarkAsClosed()
}
