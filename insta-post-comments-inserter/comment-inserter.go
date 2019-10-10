package insta_post_comments_inserter

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/codeuniversity/smag-mvp/models"
	"github.com/codeuniversity/smag-mvp/service"
	_ "github.com/lib/pq"
	"github.com/segmentio/kafka-go"
	"time"
)

type InstaCommentInserter struct {
	qReader    *kafka.Reader
	errQWriter *kafka.Writer
	db         *sql.DB
	*service.Executor
}

func New(kafkaAddress string, postgresHost string) *InstaCommentInserter {
	p := &InstaCommentInserter{}
	p.qReader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{kafkaAddress},
		GroupID:        "insta_inserter_group1",
		Topic:          "insta_comments_info",
		CommitInterval: time.Minute * 40,
	})
	p.errQWriter = kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{kafkaAddress},
		Topic:    "post_comment__insta_inserter_errors",
		Balancer: &kafka.LeastBytes{},
		Async:    false,
	})
	p.Executor = service.New()
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s user=postgres dbname=instascraper sslmode=disable", postgresHost))
	if err != nil {
		panic(err)
	}
	p.db = db
	return p
}

func (c *InstaCommentInserter) Run() {
	defer func() {
		c.MarkAsStopped()
	}()

	fmt.Println("starting Comments inserter")
	for c.IsRunning() {
		m, err := c.qReader.FetchMessage(context.Background())
		if err != nil {
			fmt.Println(err)
			break
		}
		info := &models.InstaComment{}
		err = json.Unmarshal(m.Value, info)
		if err != nil {
			panic(err)
		}
		fmt.Println("inserting: ", info.OwnerUsername)

		err = c.insertComment(info)

		if err != nil {
			panic(fmt.Errorf("comments inserter failed %s ", err))
		}
		c.qReader.CommitMessages(context.Background(), m)
	}
}

func (c *InstaCommentInserter) findOrCreateUser(username string) (userID int, err error) {
	err = c.db.QueryRow("Select id from users where user_name = $1", username).Scan(&userID)

	if err != nil {
		if err != sql.ErrNoRows {
			return 0, err
		}

		var insertedUserID int
		err := c.db.QueryRow(`INSERT INTO users(user_name) VALUES($1) RETURNING id`, username).Scan(&insertedUserID)
		if err != nil {
			return 0, err
		}

		userID = int(insertedUserID)
	}

	return userID, nil
}

func (c *InstaCommentInserter) findOrCreatePost(postId string) (postID int, err error) {
	err = c.db.QueryRow("Select id from posts where post_id = $1", postId).Scan(&postID)

	if err != nil {
		if err != sql.ErrNoRows {
			return 0, err
		}

		var insertedUserID int
		err := c.db.QueryRow(`INSERT INTO posts(post_id) VALUES($1) RETURNING id`, postId).Scan(&insertedUserID)
		if err != nil {
			return 0, err
		}

		postID = int(insertedUserID)
	}

	return postID, nil
}

func (c *InstaCommentInserter) insertComment(p *models.InstaComment) error {

	ownerUserID, err := c.findOrCreateUser(p.OwnerUsername)

	if err != nil {
		return err
	}
	postID, err := c.findOrCreatePost(p.PostId)

	if err != nil {
		return err
	}

	_, err = c.db.Exec(`INSERT INTO comments(post_id, comment_id, comment_text, owner_user_id) VALUES($1,$2,$3,$4) ON CONFLICT(comment_id) DO UPDATE SET comment_text=$3`, postID, p.Id, p.Text, ownerUserID)

	if err != nil {
		return err
	}

	return nil
}

func (c *InstaCommentInserter) Close() {
	c.Stop()
	c.WaitUntilStopped(time.Second * 3)

	c.errQWriter.Close()
	c.qReader.Close()
	c.MarkAsClosed()
}
