package utils

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

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

//WithRetries calls f up to the given `times` and returns the last error if times is reached
func WithRetries(times int, f func() error) error {
	var err error
	for i := 0; i < times; i++ {
		err = f()
		if err == nil {
			return nil
		}
		fmt.Println(err)
		time.Sleep(100 * time.Millisecond)
	}
	return err
}

// CloseOnSignal calls the closeFunc on the os signals SIGINT and SIGTERM
func CloseOnSignal(closeFunc func()) {
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		sig := <-signals
		fmt.Println("Received Signal:", sig)

		closeFunc()
	}()
}
