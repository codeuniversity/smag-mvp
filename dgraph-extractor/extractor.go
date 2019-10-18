package inserter

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgraph-io/dgo"

	"github.com/codeuniversity/smag-mvp/models"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/codeuniversity/smag-mvp/utils"

	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
)

// Extractor represents the extractor containing all clients it uses
type Extractor struct {
	qWriter  *kafka.Writer
	dgClient *dgo.Dgraph
	dgConn   *grpc.ClientConn
	*service.Executor
	currentID int
}

// New returns an initilized extractor
func New(kafkaAddress, dgraphAddress string, startID int) *Extractor {
	i := &Extractor{currentID: startID}

	i.qWriter = kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{kafkaAddress},
		Topic:    "user_follow_infos",
		Balancer: &kafka.LeastBytes{},
		Async:    true,
	})
	dg, conn := utils.GetDGraphClient(dgraphAddress)
	i.dgClient = dg
	i.dgConn = conn
	i.Executor = service.New()
	return i
}

// Run the Extractor
func (i *Extractor) Run() {
	defer func() {
		i.MarkAsStopped()
	}()

	fmt.Println("starting Extractor")
	for i.IsRunning() {
		fmt.Println("getting user with id ", i.currentID)
		info, err := getUserForID(i.dgClient, i.currentID)
		if err != nil {
			fmt.Println(err)
			i.Close()
			return
		}
		if info.CrawlTs == 0 {
			fmt.Println("skipping user with id ", i.currentID, " because he was not crawled")
		} else {
			message, err := json.Marshal(info)
			utils.PanicIfErr(err)
			fmt.Println("writing: ", i.currentID, " ", info.UserName)
			i.qWriter.WriteMessages(context.Background(), kafka.Message{Value: message})
		}
		i.currentID++
	}
}

// Close the Extractor
func (i *Extractor) Close() {
	i.Stop()
	i.WaitUntilStopped(time.Second * 3)

	i.dgConn.Close()

	i.qWriter.Close()

	i.MarkAsClosed()
}

func getUserForID(dg *dgo.Dgraph, id int) (info *models.UserFollowInfo, err error) {
	q := fmt.Sprintf(`query Me{
		me(func: uid(0x%x)){
			name
  	real_name
  	avatar_url
  	crawled_at
  	follows{
      name
    }
		}
	}`, id)

	ctx := context.Background()
	resp, err := dg.NewReadOnlyTxn().Query(ctx, q)
	if err != nil {
		return nil, err
	}
	type queryResult struct {
		Me []*models.User `json:"me"`
	}

	result := &queryResult{}
	err = json.Unmarshal(resp.GetJson(), result)

	if len(result.Me) != 1 {
		return nil, fmt.Errorf("expected 1 user to be found, but got %v", len(result.Me))
	}
	user := result.Me[0]

	info = &models.UserFollowInfo{
		UserName:  user.UserName,
		RealName:  user.RealName,
		AvatarURL: user.AvatarURL,
		Bio:       user.Bio,
		CrawlTs:   user.CrawledAt,
	}

	for _, following := range user.Follows {
		info.Followings = append(info.Followings, following.UserName)
	}
	return info, nil
}
