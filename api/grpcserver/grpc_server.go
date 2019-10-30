package grpcserver

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/url"
	"time"

	// required for postgres
	_ "github.com/lib/pq"

	"github.com/codeuniversity/smag-mvp/api/proto"
	"github.com/codeuniversity/smag-mvp/models"
	"github.com/codeuniversity/smag-mvp/utils"
	"github.com/minio/minio-go/v6"
	"google.golang.org/grpc"
)

//GrpcServer represents the gRPC Server containing the db connection and port
type GrpcServer struct {
	grpcPort string
	db       *sql.DB

	minioClient *minio.Client
	bucketName  string
	region      string
}

type scanFunc func(row *sql.Rows) (proto.User, error)

// NewGrpcServer returns initilized gRPC Server
func NewGrpcServer(grpcPort string, config models.Config) *GrpcServer {
	g := &GrpcServer{}

	g.bucketName = config.S3BucketName
	g.region = config.S3Region

	minioClient, err := minio.New(config.S3Endpoint, config.S3AccessKeyID, config.S3SecretAccessKey, config.S3UseSSL)
	utils.MustBeNil(err)
	log.Println("✅ Minio connection established")

	g.minioClient = minioClient

	postgresPassword := config.PostgresPassword
	postgresHost := config.PostgresHost

	connectionString := fmt.Sprintf("host=%s user=postgres dbname=instascraper sslmode=disable", postgresHost)
	if postgresPassword != "" {
		connectionString += " " + "password=" + postgresPassword
	}

	db, err := sql.Open("postgres", connectionString)
	utils.PanicIfNotNil(err)
	log.Println("✅ Postgres connection established")

	g.grpcPort = grpcPort
	g.db = db

	return g
}

// Listen blocks, while listening for grpc requests on the port specified in the GrpcServer struct
func (s *GrpcServer) Listen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", s.grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("✅ Start gRPC Server")

	grpcServer := grpc.NewServer()
	proto.RegisterUserSearchServiceServer(grpcServer, s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

//GetAllUsersLikeUsername returns a List of users that are like the given username
func (s *GrpcServer) GetAllUsersLikeUsername(_ context.Context, username *proto.UserNameRequest) (*proto.UserSearchResponse, error) {
	response := &proto.UserSearchResponse{}
	rows, err := s.db.Query("SELECT id, user_name, real_name, bio, avatar_url FROM users WHERE LOWER(user_name) LIKE LOWER($1)", fmt.Sprintf("%%%s%%", username.UserName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		u := proto.User{}

		err := rows.Scan(&u.Id, &u.UserName, &u.RealName, &u.Bio, &u.AvatarUrl)
		if err != nil {
			return nil, err
		}

		response.UserList = append(response.UserList, &u)
	}

	return response, nil
}

//GetUserWithUsername returns one User that equals the given username
func (s *GrpcServer) GetUserWithUsername(_ context.Context, username *proto.UserNameRequest) (*proto.User, error) {
	u := &proto.User{}
	fmt.Println(username)

	err := s.db.QueryRow("SELECT id, user_name, real_name, bio, avatar_url FROM users WHERE user_name = $1", username.UserName).Scan(&u.Id, &u.UserName, &u.RealName, &u.Bio, &u.AvatarUrl)
	if err != nil {
		return nil, err
	}

	u.Followings, err = s.getRelationsFromUser("SELECT follows.to_id as id, users.user_name FROM follows JOIN users ON follows.to_id=users.id WHERE follows.from_id=$1", u.Id, scanForIDAndUserName)
	if err != nil {
		return nil, err
	}

	u.Followers, err = s.getRelationsFromUser("SELECT follows.from_id, users.user_name FROM follows JOIN users ON follows.from_id=users.id WHERE follows.to_id=$1", u.Id, scanForIDAndUserName)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *GrpcServer) getRelationsFromUser(query string, userID string, scanFunc scanFunc) ([]*proto.User, error) {

	u := []*proto.User{}

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user := proto.User{}

		user, err = scanFunc(rows)
		if err != nil {
			return nil, err
		}

		u = append(u, &user)

	}

	return u, nil
}

//GetInstaPostsWithUserId returns all Instagram Posts of a User
func (s *GrpcServer) GetInstaPostsWithUserId(_ context.Context, request *proto.UserIdRequest) (*proto.InstaPostsResponse, error) {
	res := &proto.InstaPostsResponse{}

	rows, err := s.db.Query("SELECT id, post_id, short_code, caption, internal_picture_url FROM posts WHERE user_id=$1", request.UserId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res.UserId = request.UserId
	for rows.Next() {
		post := proto.InstaPost{}

		rows.Scan(&post.Id, &post.PostId, &post.ShortCode, &post.Caption, &post.ImgUrl)

		if post.ImgUrl != "" {
			post.ImgUrl, err = s.getURLForPost(post.ImgUrl)
			if err != nil {
				return nil, err
			}
		}

		res.InstaPosts = append(res.InstaPosts, &post)
	}

	return res, nil
}

//scanForIdAndUserName scans a sql row for user id and username
func scanForIDAndUserName(row *sql.Rows) (proto.User, error) {
	user := proto.User{}
	err := row.Scan(&user.Id, &user.UserName)
	if err != nil {
		//return nil, nil
	}

	return user, nil
}

func (s *GrpcServer) getURLForPost(object string) (string, error) {

	// Set request parameters for content-disposition.
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", fmt.Sprintf("attachment; filename=\"%s.jpg\"", object))

	// Generates a presigned url which expires in a day.
	presignedURL, err := s.minioClient.PresignedGetObject(s.bucketName, object, time.Second*24*60*60, reqParams)
	if err != nil {
		log.Println(err)
		return "", err
	}
	fmt.Println("Successfully generated presigned URL", presignedURL)
	return presignedURL.String(), nil

}
