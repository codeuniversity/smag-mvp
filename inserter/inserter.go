package inserter

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"

	"github.com/alexmorten/instascraper/models"
	"github.com/alexmorten/instascraper/utils"
	"github.com/segmentio/kafka-go"
)

// Inserter represents the scraper containing all clients it uses
type Inserter struct {
	qReader     *kafka.Reader
	qWriter     *kafka.Writer
	stopChan    chan struct{}
	stoppedChan chan struct{}
	closedChan  chan struct{}
}

// New returns an initilized scraper
func New() *Inserter {
	s := &Inserter{}
	s.qReader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{"localhost:9092"},
		GroupID:        "user_follow_inserter",
		Topic:          "user_follow_infos",
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		CommitInterval: time.Second,
	})
	s.qWriter = kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "user_names",
		Balancer: &kafka.LeastBytes{},
		Async:    true,
	})
	s.stopChan = make(chan struct{}, 1)
	s.stoppedChan = make(chan struct{}, 1)
	s.closedChan = make(chan struct{}, 1)
	return s
}

// Run the inserter
func (i *Inserter) Run() {
	defer func() {
		i.stoppedChan <- struct{}{}
	}()

	fmt.Println("starting inserter")
	for len(i.stopChan) == 0 {
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
		fmt.Println("inserted: ", info.UserName)
		i.qReader.CommitMessages(context.Background(), m)
		fmt.Println("commited: ", info.UserName)
	}
}

// Close the inserter
func (i *Inserter) Close() {
	i.stopChan <- struct{}{}
	t := time.NewTimer(time.Second * 3)
	select {
	case <-t.C:
		break
	case <-i.stoppedChan:
		t.Stop()
		break
	}

	i.qReader.Close()
	i.qWriter.Close()
	i.closedChan <- struct{}{}
}

// WaitUntilClosed waits until the Close func call of the inserter is finished
func (i *Inserter) WaitUntilClosed() {
	<-i.closedChan
}

// InsertUserFollowInfo inserts the user follow info into dgraph, while writting userNames that don't exist in the graph yet
// into the specified kafka topic
func (i *Inserter) InsertUserFollowInfo(followInfo *models.UserFollowInfo) {
	for _, follower := range followInfo.Followers {
		p := &models.User{
			Name: follower,
			Follows: []*models.User{
				{Name: followInfo.UserName},
			},
		}

		i.insertUser(p)
	}

	p := &models.User{
		Name: followInfo.UserName,
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
	dg, conn := utils.GetDGraphClient()
	defer conn.Close()

	uid, created := getOrCreateUIDForUserWithRetries(dg, p.Name)
	i.handleCreatedUser(p.Name, uid, created)
	p.UID = uid
	for _, followed := range p.Follows {
		uid, created := getOrCreateUIDForUserWithRetries(dg, followed.Name)
		i.handleCreatedUser(followed.Name, uid, created)
		followed.UID = uid
	}

	mu := &api.Mutation{
		CommitNow: true,
	}

	pb, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	mu.SetJson = pb

	err = utils.WithRetries(5, func() error {
		ctx := context.Background()
		_, err = dg.NewTxn().Mutate(ctx, mu)
		return err
	})
	if err != nil {
		panic(err)
	}
}

func (i *Inserter) handleCreatedUser(userName, uid string, created bool) {
	if created {
		i.qWriter.WriteMessages(context.Background(), kafka.Message{
			Value: []byte(userName),
		})
	}
}

func getOrCreateUIDForUser(dg *dgo.Dgraph, name string) (uid string, created bool, err error) {
	q := `query Me($name: string){
		me(func: eq(name, $name)){
			uid
		}
	}`
	ctx := context.Background()
	resp, err := dg.NewReadOnlyTxn().QueryWithVars(ctx, q, map[string]string{"$name": name})
	if err != nil {
		return "", false, err
	}
	type queryResult struct {
		Me []*models.User `json:"me"`
	}
	result := &queryResult{}
	err = json.Unmarshal(resp.GetJson(), result)
	if err != nil {
		return "", false, err
	}
	if len(result.Me) > 0 {
		return result.Me[0].UID, false, nil
	}
	mu := &api.Mutation{
		CommitNow: true,
	}
	p := &models.User{Name: name}
	pb, err := json.Marshal(p)
	if err != nil {
		return "", false, err
	}

	ctx = context.Background()
	mu.SetJson = pb
	assigned, err := dg.NewTxn().Mutate(ctx, mu)
	if err != nil {
		return "", false, err
	}

	// Assigned uids for nodes which were created would be returned in the assigned.Uids map.
	return assigned.Uids["blank-0"], true, nil
}

func getOrCreateUIDForUserWithRetries(dg *dgo.Dgraph, name string) (uid string, created bool) {
	err := utils.WithRetries(5, func() error {
		var err error
		uid, created, err = getOrCreateUIDForUser(dg, name)
		return err
	})
	if err != nil {
		panic(err)
	}
	return
}
