package utils

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"os"
	"strconv"
	"time"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"google.golang.org/grpc"
)

//GetDGraphClient retuns an intialized dgrapg client and a connection that should be closed once the client is discarded.
// panics if dgraog can not be connected to
func GetDGraphClient(dgraphAddress string) (*dgo.Dgraph, *grpc.ClientConn) {
	conn, err := grpc.Dial(dgraphAddress, grpc.WithInsecure())
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

// if qWriter is nil, user discovery is disabled
func HandleCreatedUser(qWriter *kafka.Writer, userName string) {
	if qWriter != nil {
		qWriter.WriteMessages(context.Background(), kafka.Message{
			Value: []byte(userName),
		})
	}
}

//GetStringFromEnvWithDefault returns default Value if OS Environment Variable is not set
func GetStringFromEnvWithDefault(enVarName, defaultValue string) string {
	envValue := os.Getenv(enVarName)
	if envValue == "" {
		return defaultValue
	}

	return envValue
}

//MustGetStringFromEnv panics if OS Environment Variable is not set
func MustGetStringFromEnv(enVarName string) string {
	envValue := os.Getenv(enVarName)
	if envValue == "" {
		panic(fmt.Sprintf("%s must not be empty", enVarName))
	}

	return envValue
}

// GetBoolFromEnvWithDefault parses an OS Environment Variable as bool
func GetBoolFromEnvWithDefault(enVarName string, defaultValue bool) bool {
	envValue := os.Getenv(enVarName)
	if envValue == "" {
		return defaultValue
	}

	envBool, err := strconv.ParseBool(envValue)
	if err != nil {
		panic(fmt.Errorf("couldn't parse %s as bool: %s", enVarName, err))
	}

	return envBool
}
