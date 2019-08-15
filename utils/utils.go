package utils

import (
	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"google.golang.org/grpc"
)

//GetDGraphClient retuns an intialized dgrapg client and a connection that should be closed once the client is discarded.
// panics if dgraog can not be connected to
func GetDGraphClient() (*dgo.Dgraph, *grpc.ClientConn) {
	conn, err := grpc.Dial("127.0.0.1:9080", grpc.WithInsecure())
	if err != nil {
		panic("While trying to dial gRPC")
	}

	dc := api.NewDgraphClient(conn)
	dg := dgo.NewDgraphClient(dc)
	return dg, conn
}
