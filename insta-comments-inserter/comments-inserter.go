package insta_comments_inserter

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/codeuniversity/smag-mvp/models"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/codeuniversity/smag-mvp/utils"
	_ "github.com/lib/pq"
	"github.com/segmentio/kafka-go"
	"time"
)

type InstaCommentInserter struct {
	commentsQReader *kafka.Reader
	userQWriter     *kafka.Writer
	db              *sql.DB
	*service.Executor
}

func New(postgresHost, postgresPassword string, commentsQReader *kafka.Reader, userQWriter *kafka.Writer) *InstaCommentInserter {
	p := &InstaCommentInserter{}
	p.commentsQReader = commentsQReader
	p.userQWriter = userQWriter
	p.Executor = service.New()

	connectionString := fmt.Sprintf("host=%s user=postgres dbname=instascraper sslmode=disable", postgresHost)
	if postgresPassword != "" {
		connectionString += " " + "password=" + postgresPassword
	}

	db, err := sql.Open("postgres", connectionString)
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
		m, err := c.commentsQReader.FetchMessage(context.Background())
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
		c.commentsQReader.CommitMessages(context.Background(), m)
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

		utils.HandleCreatedUser(c.userQWriter, username)
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

	c.userQWriter.Close()
	c.commentsQReader.Close()
	c.MarkAsClosed()
}
