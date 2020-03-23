package inserter

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/codeuniversity/smag-mvp/insta/models"
	"github.com/codeuniversity/smag-mvp/worker"

	// necessary for "database/sql"
	_ "github.com/lib/pq"
	"github.com/segmentio/kafka-go"
)

// InstaPostInserter is responsible for inserting models.InstagramPost
type InstaPostInserter struct {
	*worker.Worker

	postQReader *kafka.Reader

	db *sql.DB
}

// New retuns an initialized InstaPostInserter
func New(postgresHost, postgresPassword string, postQReader *kafka.Reader) *InstaPostInserter {
	i := &InstaPostInserter{}
	i.postQReader = postQReader

	connectionString := fmt.Sprintf("host=%s user=postgres dbname=instascraper sslmode=disable", postgresHost)
	if postgresPassword != "" {
		connectionString += " " + "password=" + postgresPassword
	}

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}
	i.db = db

	i.Worker = worker.Builder{}.WithName("insta_posts_inserter").
		WithWorkStep(i.runStep).
		WithStopTimeout(10*time.Second).
		AddShutdownHook("postQReader", postQReader.Close).
		AddShutdownHook("postgres_connection", db.Close).
		MustBuild()

	return i
}

func (i *InstaPostInserter) runStep() error {
	message, err := i.postQReader.FetchMessage(context.Background())
	if err != nil {
		return err
	}

	var post models.InstagramPost
	err = json.Unmarshal(message.Value, &post)
	if err != nil {
		return err
	}
	log.Println(post)
	post.Caption = strings.ReplaceAll(post.Caption, "\u0000", "")

	postID, err := i.insertPost(post)

	if err != nil {
		return fmt.Errorf("posts inserter insertPost() failed %s ", err)
	}

	err = i.insertTaggedUser(postID, post.TaggedUsers)

	if err != nil {
		return fmt.Errorf("posts inserter insertTaggedUser() failed %s ", err)
	}
	log.Println("Insert Post: ", post.ShortCode)
	return i.postQReader.CommitMessages(context.Background(), message)
}

func (i *InstaPostInserter) findOrCreateUser(username string) (userID int, err error) {
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

func (i *InstaPostInserter) insertTaggedUser(postID int, taggedUsers []string) error {
	if taggedUsers == nil {
		return nil
	}

	for _, username := range taggedUsers {

		userID, err := i.findOrCreateUser(username)

		if err != nil {
			return err
		}
		var taggedID int
		err = i.db.QueryRow("Select id from post_tagged_users where post_id=$1 AND user_id= $2", postID, userID).Scan(&taggedID)
		if err == sql.ErrNoRows {
			_, err = i.db.Exec("Insert INTO post_tagged_users(post_id,user_id) VALUES($1,$2)", postID, userID)
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	}
	return nil
}

func (i *InstaPostInserter) insertPost(post models.InstagramPost) (int, error) {
	userID, err := i.findOrCreateUser(post.UserName)
	if err != nil {
		return 0, err
	}

	var postID int
	err = i.db.QueryRow(`INSERT INTO posts(user_id, post_id, short_code, picture_url, caption) VALUES($1,$2,$3,$4,$5) ON CONFLICT(post_id) DO UPDATE SET short_code=$3, picture_url=$4, caption=$5 RETURNING id`, userID, post.PostID, post.ShortCode, post.PictureURL, post.Caption).Scan(&postID)
	if err != nil {
		return 0, err
	}

	return postID, nil
}
