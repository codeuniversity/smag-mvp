package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gocolly/colly"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"

	"github.com/alexmorten/instascraper/utils"
)

type user struct {
	Followers  []string `json:"followers"`
	Followings []string `json:"followings"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a username as the first argument")
		return
	}

	userNameArg := os.Args[1]
	// prepareDB()
	followers, err := getUserNamesIn(fmt.Sprintf("http://picdeer.com/%s/followers", userNameArg))
	handleErr(err)
	followings, err := getUserNamesIn(fmt.Sprintf("http://picdeer.com/%s/followings", userNameArg))
	handleErr(err)
	result, err := json.Marshal(&user{Followers: followers, Followings: followings})
	handleErr(err)
	fmt.Println(string(result))

	for _, follower := range followers {
		p := &Person{
			Name: follower,
			Followings: []*Person{
				{Name: userNameArg},
			},
		}

		insertUser(p)
	}

	for _, followings := range followings {
		p := &Person{
			Name: followings,
			Followers: []*Person{
				{Name: userNameArg},
			},
		}

		insertUser(p)
	}

	p := &Person{
		Name: userNameArg,
	}

	for _, following := range followings {
		p.Followings = append(p.Followings, &Person{
			Name: following,
		})
	}

	for _, follower := range followers {
		p.Followers = append(p.Followers, &Person{
			Name: follower,
		})
	}

	insertUser(p)
}

func getUserNamesIn(url string) ([]string, error) {
	userNames := []string{}

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		// fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		// log.Println("Something went wrong:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		// fmt.Println("Visited", r.Request.URL)
	})

	c.OnHTML("p.grid-user-identifier-1 a", func(e *colly.HTMLElement) {
		url := e.Attr("href")
		parts := strings.Split(url, "/")
		if len(parts) >= 2 {
			name := parts[len(parts)-1]
			userNames = append(userNames, name)
		}
	})

	c.OnScraped(func(r *colly.Response) {
		// fmt.Println("Finished", r.Request.URL)
	})

	err := c.Visit(url)
	if err != nil {
		return nil, err
	}

	return userNames, nil
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

// Person is the struct used for serializing users to dgraph
type Person struct {
	UID        string    `json:"uid,omitempty"`
	Name       string    `json:"name,omitempty"`
	Followers  []*Person `json:"follows,omitempty"`
	Followings []*Person `json:"followed,omitempty"`
}

func insertUser(p *Person) {
	dg, conn := utils.GetDGraphClient()
	defer conn.Close()

	p.UID = getOrCreateUIDForUser(dg, p.Name)
	for _, followed := range p.Followings {
		followed.UID = getOrCreateUIDForUser(dg, followed.Name)
	}
	for _, following := range p.Followers {
		following.UID = getOrCreateUIDForUser(dg, following.Name)
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

func getOrCreateUIDForUser(dg *dgo.Dgraph, name string) string {
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
		Me []*Person `json:"me"`
	}
	result := &queryResult{}
	err = json.Unmarshal(resp.GetJson(), result)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
	if len(result.Me) > 0 {
		fmt.Println("returned id")
		return result.Me[0].UID
	}
	fmt.Println("create id")
	mu := &api.Mutation{
		CommitNow: true,
	}
	p := &Person{Name: name}
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
	return assigned.Uids["blank-0"]
}
