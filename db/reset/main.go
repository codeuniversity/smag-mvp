package main

import (
	"os"
	"context"

	"github.com/alexmorten/instascraper/utils"
	"github.com/dgraph-io/dgo/protos/api"
)

func main() {
	prepareDB()
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func prepareDB() {
	dgraphAddress := os.Getenv("DGRAPH_ADDRESS")
	if dgraphAddress == "" {
		dgraphAddress = "127.0.0.1:9080"
	}
	dg, conn := utils.GetDGraphClient(dgraphAddress)
	defer conn.Close()

	op := &api.Operation{DropAll: true}

	ctx := context.Background()
	handleErr(dg.Alter(ctx, op))

	op = &api.Operation{}
	op.Schema = `
	name: string @index(exact) @upsert .
	real_name: string @index(fulltext, term) .
	follows: uid @count @reverse .
	crawled_at: int .
	bio: string @index(fulltext, term) .
	`

	ctx = context.Background()
	handleErr(dg.Alter(ctx, op))

}
