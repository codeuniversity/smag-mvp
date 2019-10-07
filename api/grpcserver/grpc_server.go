package grpcserver

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"

	// required for postgres
	_ "github.com/lib/pq"

	"github.com/codeuniversity/smag-mvp/api/proto"
	"google.golang.org/grpc"
)

//GrpcServer holds port of server
type GrpcServer struct {
	Users    *proto.User
	grpcPort int
	db       *sql.DB
}

// NewGrpcServer creates a GrpcServer
func NewGrpcServer(grpcPort int, postgresHost string) *GrpcServer {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s user=postgres dbname=instascraper sslmode=disable", postgresHost))
	if err != nil {
		panic(err)
	}

	return &GrpcServer{
		grpcPort: grpcPort,
		db:       db,
	}
}

// Listen blocks, while listening for grpc requests on the port specified in the GrpcServer struct
func (s *GrpcServer) Listen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", s.grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	proto.RegisterUserSearchServiceServer(grpcServer, s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

//GetUserWithUsername searches for User in Database
func (s *GrpcServer) GetUserWithUsername(_ context.Context, username *proto.UserName) (*proto.UserSearchResponse, error) {
	fmt.Println("Hallo")
	query := fmt.Sprintf("SELECT * FROM users WHERE user_name LIKE '%s'", username.UserName)
	response := &proto.UserSearchResponse{}

	rows, err := s.db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		u := proto.User{}

		rows.Scan(&u.UserName, &u.RealName, &u.Bio, &u.AvatarUrl, &u.FollowingsCount, &u.FollowersCount)

		// TODO: error handling
		response.UserList = append(response.UserList, &u)
	}
	if err != nil {
		log.Fatal(err)
	}

	return response, nil
}
