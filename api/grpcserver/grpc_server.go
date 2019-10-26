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
	"github.com/codeuniversity/smag-mvp/utils"
	"google.golang.org/grpc"
)

//GrpcServer represents the gRPC Server containing the db connection and port
type GrpcServer struct {
	grpcPort string
	db       *sql.DB
}

// NewGrpcServer returns initilized gRPC Server
func NewGrpcServer(postgresHost string, postgresPassword string, grpcPort string) *GrpcServer {

	connectionString := fmt.Sprintf("host=%s user=postgres dbname=instascraper sslmode=disable", postgresHost)
	if postgresPassword != "" {
		connectionString += " " + "password=" + postgresPassword
	}

	db, err := sql.Open("postgres", connectionString)
	utils.PanicIfNotNil(err)

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

	fmt.Println("Start gRPC Server")

	grpcServer := grpc.NewServer()
	proto.RegisterUserSearchServiceServer(grpcServer, s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

//GetAllUsersLikeUsername returns a List of users that are like the given username
func (s *GrpcServer) GetAllUsersLikeUsername(_ context.Context, username *proto.UserSearchRequest) (*proto.UserSearchResponse, error) {
	response := &proto.UserSearchResponse{}
	rows, err := s.db.Query("SELECT id, user_name, real_name, bio, avatar_url FROM users WHERE LOWER(user_name) LIKE LOWER($1)", fmt.Sprintf("%%%s%%", username.UserName))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		u := proto.User{}

		rows.Scan(&u.Id, &u.UserName, &u.RealName, &u.Bio, &u.AvatarUrl)

		// TODO: error handling
		response.UserList = append(response.UserList, &u)
	}
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(response)
	return response, nil
}

//GetUserWithUsername returns one User that equals the given username
func (s *GrpcServer) GetUserWithUsername(_ context.Context, username *proto.UserSearchRequest) (*proto.User, error) {
	u := &proto.User{}
	fmt.Println(username)

	err := s.db.QueryRow("SELECT id, user_name, real_name, bio, avatar_url FROM users WHERE user_name = $1", username.UserName).Scan(&u.Id, &u.UserName, &u.RealName, &u.Bio, &u.AvatarUrl)
	if err != nil {
		return nil, err
	}

	u.Followings, err = s.getRelationsFromUser("SELECT follows.to_id as id, users.user_name FROM follows JOIN users ON follows.to_id=users.id WHERE follows.from_id=$1", u.Id)
	if err != nil {
		return nil, err
	}

	u.Followers, err = s.getRelationsFromUser("SELECT follows.from_id, users.user_name FROM follows JOIN users ON follows.from_id=users.id WHERE follows.to_id=$1", u.Id)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *GrpcServer) getRelationsFromUser(query string, userID string) ([]*proto.User, error) {

	u := []*proto.User{}

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user := proto.User{}

		rows.Scan(&user.Id, &user.UserName)

		u = append(u, &user)

	}

	return u, nil
}
