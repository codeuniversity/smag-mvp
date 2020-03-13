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

// InstaLikesInserter inserts likes into postgres
type InstaLikesInserter struct {
	*worker.Worker

	likesQReader *kafka.Reader

	db *sql.DB
}

// New returns an initialized InstaLikesInserter
func New(postgresHost, postgresPassword string, likesQReader *kafka.Reader) *InstaLikesInserter {
	i := &InstaLikesInserter{}
	i.likesQReader = likesQReader

	connectionString := fmt.Sprintf("host=%s user=postgres dbname=instascraper sslmode=disable", postgresHost)
	if postgresPassword != "" {
		connectionString += " " + "password=" + postgresPassword
	}

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}
	i.db = db

	b := worker.Builder{}.WithName("insta_likes_inserter").
		WithWorkStep(i.runStep).
		WithStopTimeout(10*time.Second).
		AddShutdownHook("likesQReader", likesQReader.Close).
		AddShutdownHook("postgres_client", db.Close)

	i.Worker = b.MustBuild()

	return i
}

func (i *InstaLikesInserter) runStep() error {
	m, err := i.likesQReader.FetchMessage(context.Background())
	if err != nil {
		return err
	}
	info := &models.InstaLike{}
	err = json.Unmarshal(m.Value, info)
	if err != nil {
		return err
	}
	log.Println("inserting: ", info.OwnerUsername)

	err = i.insertLike(info)

	if err != nil {
		return fmt.Errorf("likes inserter failed %s ", err)
	}
	return i.likesQReader.CommitMessages(context.Background(), m)
}

func (i *InstaLikesInserter) findOrCreateUser(username string) (userID int, err error) {
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

func (i *InstaLikesInserter) findOrCreatePost(externalPostID string) (postID int, err error) {
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

func (i *InstaLikesInserter) insertLike(p *models.InstaLike) error {

	ownerUserID, err := i.findOrCreateUser(p.OwnerUsername)

	if err != nil {
		return err
	}

	postID, err := i.findOrCreatePost(p.PostID)

	if err != nil {
		return err
	}

	_, err = i.db.Exec(`INSERT INTO post_likes(like_id, user_id, post_id) VALUES($1,$2,$3) ON CONFLICT(like_id) DO NOTHING`, p.ID, ownerUserID, postID)

	if err != nil {
		return err
	}

	return nil
}
