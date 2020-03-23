package inserter

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/codeuniversity/smag-mvp/insta/models"
	"github.com/codeuniversity/smag-mvp/worker"

	// necessary for "database/sql"
	_ "github.com/lib/pq"
	"github.com/segmentio/kafka-go"
)

// InstaCommentInserter inserts comments into postgres
type InstaCommentInserter struct {
	*worker.Worker

	commentsQReader *kafka.Reader

	db *sql.DB
}

// New returns an initialized InstaCommentInserter
func New(postgresHost, postgresPassword string, commentsQReader *kafka.Reader) *InstaCommentInserter {
	i := &InstaCommentInserter{}
	i.commentsQReader = commentsQReader

	connectionString := fmt.Sprintf("host=%s user=postgres dbname=instascraper sslmode=disable", postgresHost)
	if postgresPassword != "" {
		connectionString += " " + "password=" + postgresPassword
	}

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}
	i.db = db

	b := worker.Builder{}.WithName("insta_comments_inserter").
		WithWorkStep(i.runStep).
		WithStopTimeout(10*time.Second).
		AddShutdownHook("commentsQReader", commentsQReader.Close).
		AddShutdownHook("postgres_client", db.Close)

	i.Worker = b.MustBuild()

	return i
}

func (i *InstaCommentInserter) runStep() error {
	m, err := i.commentsQReader.FetchMessage(context.Background())
	if err != nil {
		return err
	}
	info := &models.InstaComment{}
	err = json.Unmarshal(m.Value, info)
	if err != nil {
		return err
	}
	log.Println("inserting: ", info.OwnerUsername)

	err = i.insertComment(info)

	if err != nil {
		return fmt.Errorf("comments inserter failed %s ", err)
	}
	return i.commentsQReader.CommitMessages(context.Background(), m)
}

func (i *InstaCommentInserter) findOrCreateUser(username string) (userID int, err error) {
	err = i.db.QueryRow("Select id from users where user_name = $1", username).Scan(&userID)

	if err != nil {
		if err != sql.ErrNoRows {
			return 0, err
		}

		var insertedUserID int
		err := i.db.QueryRow(`INSERT INTO users(user_name) VALUES($1) RETURNING id`, username).Scan(&insertedUserID)
		if err != nil {
			return 0, err
		}

		userID = int(insertedUserID)
	}

	return userID, nil
}

func (i *InstaCommentInserter) findOrCreatePost(externalPostID string) (postID int, err error) {
	err = i.db.QueryRow("Select id from posts where post_id = $1", externalPostID).Scan(&postID)

	if err != nil {
		if err != sql.ErrNoRows {
			return 0, err
		}

		var insertedUserID int
		err := i.db.QueryRow(`INSERT INTO posts(post_id) VALUES($1) RETURNING id`, externalPostID).Scan(&insertedUserID)
		if err != nil {
			return 0, err
		}

		postID = int(insertedUserID)
	}

	return postID, nil
}

func (i *InstaCommentInserter) insertComment(p *models.InstaComment) error {

	ownerUserID, err := i.findOrCreateUser(p.OwnerUsername)

	if err != nil {
		return err
	}
	postID, err := i.findOrCreatePost(p.PostID)

	if err != nil {
		return err
	}

	_, err = i.db.Exec(`INSERT INTO comments(post_id, comment_id, comment_text, owner_user_id) VALUES($1,$2,$3,$4) ON CONFLICT(comment_id) DO UPDATE SET comment_text=$3`, postID, p.ID, p.Text, ownerUserID)

	if err != nil {
		return err
	}

	return nil
}
