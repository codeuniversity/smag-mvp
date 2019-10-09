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

//GetAllUsersLikeUsername returns a List of users that are like the given username
func (s *GrpcServer) GetAllUsersLikeUsername(_ context.Context, username *proto.UserName) (*proto.UserSearchResponse, error) {
	response := &proto.UserSearchResponse{}
	rows, err := s.db.Query("SELECT user_nanme, real_name, bio, avatar_url WHERE LOWER(user_name) LIKE LOWER($1)", username.UserName)
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

//GetUserWithUsername returns one User that equals the given username
func (s *GrpcServer) GetUserWithUsername(_ context.Context, username *proto.UserName) (*proto.User, error) {
	u := &proto.User{}

	row := s.db.QueryRow("SELECT user_nanme, real_name, bio, avatar_url FROM users WHERE user_name = $1", username.UserName)

	row.Scan(&u.UserName, &u.RealName, &u.Bio)
	return u, nil
}
