package inserter

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"

	"github.com/codeuniversity/smag-mvp/models"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/codeuniversity/smag-mvp/utils"

	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
)

// Inserter represents the scraper containing all clients it uses
type Inserter struct {
	qReader  *kafka.Reader
	qWriter  *kafka.Writer
	dgClient *dgo.Dgraph
	dgConn   *grpc.ClientConn
	*service.Executor
}

// New returns an initilized scraper
func New(dgraphAddress string, qReader *kafka.Reader, qWriter *kafka.Writer) *Inserter {
	i := &Inserter{}
	i.qReader = qReader
	i.qWriter = qWriter

	dg, conn := utils.GetDGraphClient(dgraphAddress)
	i.dgClient = dg
	i.dgConn = conn
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

	i.dgConn.Close()
	i.qReader.Close()
	if i.qWriter != nil {
		i.qWriter.Close()
	}

	i.MarkAsClosed()
}

// InsertUserFollowInfo inserts the user follow info into dgraph, while writting userNames that don't exist in the graph yet
// into the specified kafka topic
func (i *Inserter) InsertUserFollowInfo(followInfo *models.UserFollowInfo) {
	p := &models.User{
		UserName:  followInfo.UserName,
		RealName:  followInfo.RealName,
		AvatarURL: followInfo.AvatarURL,
		Bio:       followInfo.Bio,
		CrawledAt: followInfo.CrawlTs,
	}

	for _, following := range followInfo.Followings {
		p.Follows = append(p.Follows, &models.User{
			UserName: following,
		})
	}

	i.insertUser(p)
}

func (i *Inserter) insertUser(p *models.User) {
	uid, created := getOrCreateUIDForUserWithRetries(i.dgClient, p.UserName)
	i.handleCreatedUser(p.UserName, uid, created)
	p.UID = uid
	for _, followed := range p.Follows {
		uid, created := getOrCreateUIDForUserWithRetries(i.dgClient, followed.UserName)
		i.handleCreatedUser(followed.UserName, uid, created)
		followed.UID = uid
	}

	mu := &api.Mutation{
		CommitNow: true,
	}

	pb, err := json.Marshal(p)
	utils.PanicIfErr(err)
	mu.SetJson = pb

	err = utils.WithRetries(5, func() error {
		ctx := context.Background()
		_, err = i.dgClient.NewTxn().Mutate(ctx, mu)
		return err
	})
	utils.PanicIfErr(err)
}

func (i *Inserter) handleCreatedUser(userName, uid string, created bool) {
	// if qWriter is nil, user discovery is disabled
	if created && i.qWriter != nil {
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
	p := &models.User{UserName: name}
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
	utils.PanicIfErr(err)
	return
}
