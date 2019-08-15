package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"

	"github.com/alexmorten/instascraper/models"
	"github.com/alexmorten/instascraper/utils"
)

// InsertUserFollowInfo inserts the user follow info into dgraph, while writting userNames that don't exist in the graph yet
// into the specified kafka topic
func InsertUserFollowInfo(followInfo *models.UserFollowInfo) {
	for _, follower := range followInfo.Followers {
		p := &models.User{
			Name: follower,
			Followings: []*models.User{
				{Name: followInfo.UserName},
			},
		}

		insertUser(p)
	}

	for _, followings := range followInfo.Followings {
		p := &models.User{
			Name: followings,
			Followers: []*models.User{
				{Name: followInfo.UserName},
			},
		}

		insertUser(p)
	}

	p := &models.User{
		Name: followInfo.UserName,
	}

	for _, following := range followInfo.Followings {
		p.Followings = append(p.Followings, &models.User{
			Name: following,
		})
	}

	for _, follower := range followInfo.Followers {
		p.Followers = append(p.Followers, &models.User{
			Name: follower,
		})
	}

	insertUser(p)
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func insertUser(p *models.User) {
	dg, conn := utils.GetDGraphClient()
	defer conn.Close()

	uid, created := getOrCreateUIDForUser(dg, p.Name)
	handleCreatedUser(p.Name, uid, created)
	p.UID = uid
	for _, followed := range p.Followings {
		uid, created := getOrCreateUIDForUser(dg, followed.Name)
		handleCreatedUser(followed.Name, uid, created)
		followed.UID = uid
	}
	for _, following := range p.Followers {
		uid, created := getOrCreateUIDForUser(dg, following.Name)
		handleCreatedUser(following.Name, uid, created)
		following.UID = uid
	}

	mu := &api.Mutation{
		CommitNow: true,
	}

	pb, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	mu.SetJson = pb
	assigned, err := dg.NewTxn().Mutate(ctx, mu)
	if err != nil {
		panic(err)
	}
	fmt.Println(assigned)
}

func handleCreatedUser(userName, uid string, created bool) {
	if created {
		fmt.Println("encountered new user:", userName, "(", uid, ")")
	}
}

func getOrCreateUIDForUser(dg *dgo.Dgraph, name string) (uid string, created bool) {
	q := `query Me($name: string){
		me(func: eq(name, $name)){
			uid
		}
	}`
	ctx := context.Background()
	resp, err := dg.NewReadOnlyTxn().QueryWithVars(ctx, q, map[string]string{"$name": name})
	if err != nil {
		panic(err)
	}
	fmt.Println(string(resp.GetJson()))
	type queryResult struct {
		Me []*models.User `json:"me"`
	}
	result := &queryResult{}
	err = json.Unmarshal(resp.GetJson(), result)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
	if len(result.Me) > 0 {
		fmt.Println("returned id")
		return result.Me[0].UID, false
	}
	fmt.Println("created id")
	mu := &api.Mutation{
		CommitNow: true,
	}
	p := &models.User{Name: name}
	pb, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}

	ctx = context.Background()
	mu.SetJson = pb
	assigned, err := dg.NewTxn().Mutate(ctx, mu)
	if err != nil {
		panic(err)
	}

	// Assigned uids for nodes which were created would be returned in the assigned.Uids map.
	return assigned.Uids["blank-0"], true
}
