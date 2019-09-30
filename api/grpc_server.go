package usersearch

import (
	"fmt"
	"log"
	"net"
	"github.com/codeuniversity/mvp-smag/api/proto"

	"google.golang.org/grpc"
)

type GrpcServer struct {
	grpcPort int
}

// NewGrpcServer creates a GrpcServer
func NewGrpcServer(grpcPort int) *GrpcServer {
	return &GrpcServer{
		grpcPort: grpcPort,
	}
}

// Listen blocks, while listening for grpc requests on the port specified in the GrpcServer struct
func (s *GrpcServer) Listen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", s.grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	proto.

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
