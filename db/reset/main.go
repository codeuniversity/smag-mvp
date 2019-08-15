package main

import (
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
	dg, conn := utils.GetDGraphClient()
	defer conn.Close()

	op := &api.Operation{DropAll: true}

	ctx := context.Background()
	handleErr(dg.Alter(ctx, op))

	op = &api.Operation{}
	op.Schema = `
	name: string @index(exact) @upsert .`

	ctx = context.Background()
	handleErr(dg.Alter(ctx, op))

}
