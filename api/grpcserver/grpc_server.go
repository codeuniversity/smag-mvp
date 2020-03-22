package grpcserver

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	// required for postgres
	_ "github.com/lib/pq"
	"github.com/segmentio/kafka-go"

	"github.com/codeuniversity/smag-mvp/api/proto"
	"github.com/minio/minio-go/v6"
	"google.golang.org/grpc"

	"github.com/codeuniversity/smag-mvp/config"
	"github.com/codeuniversity/smag-mvp/utils"

	"github.com/codeuniversity/smag-mvp/elastic"
	"github.com/codeuniversity/smag-mvp/elastic/search/faces"
	recognition "github.com/codeuniversity/smag-mvp/faces/proto"
	analyzer "github.com/codeuniversity/smag-mvp/nlp/frequency-analyzer"
)

const (
	signedURLExpirationTime = time.Minute * 60
)

//GrpcServer represents the gRPC Server containing the db connection and port
type GrpcServer struct {
	grpcPort string
	db       *sql.DB

	minioClient        *minio.Client
	downloadBucketName string
	region             string
	imageUploadBucket  string
	userNamesWriter    *kafka.Writer

	facesClient *faces.Client
	nlpAnalyzer *analyzer.Analyzer
	citiesMap   map[string][]string
}

type scanFunc func(row *sql.Rows) (proto.User, error)

// NewGrpcServer returns initilized gRPC Server
func NewGrpcServer(grpcPort string, userNamesWriter *kafka.Writer, s3Config *config.S3Config, imageUploadBucket string, postgresConfig *config.PostgresConfig, esHosts []string, recognitionServiceAddress string) *GrpcServer {
	s := &GrpcServer{}

	// TODO: load cities.json
	jsonFile, err := os.Open("nlp/frequency-analyzer/cities.json")
	defer jsonFile.Close()
	if err != nil {
		log.Fatalln(err)
	}
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err := json.Unmarshal(byteValue, &s.citiesMap); err != nil {
		panic(err)
	}
	s.userNamesWriter = userNamesWriter

	s.downloadBucketName = s3Config.S3BucketName
	s.region = s3Config.S3Region
	s.imageUploadBucket = imageUploadBucket

	minioClient, err := minio.New(s3Config.S3Endpoint, s3Config.S3AccessKeyID, s3Config.S3SecretAccessKey, s3Config.S3UseSSL)
	utils.MustBeNil(err)
	log.Println("✅ Minio connection established")
	s.nlpAnalyzer = analyzer.New(esHosts)

	s.minioClient = minioClient

	postgresPassword := postgresConfig.PostgresPassword
	postgresHost := postgresConfig.PostgresHost

	connectionString := fmt.Sprintf("host=%s user=postgres dbname=instascraper sslmode=disable", postgresHost)
	if postgresPassword != "" {
		connectionString += " " + "password=" + postgresPassword
	}

	db, err := sql.Open("postgres", connectionString)
	utils.PanicIfNotNil(err)
	log.Println("✅ Postgres connection established")

	con, err := grpc.Dial(recognitionServiceAddress, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	faceReconClient := recognition.NewFaceRecognizerClient(con)
	esClient := elastic.InitializeElasticSearch(esHosts)
	s.facesClient = &faces.Client{
		FaceRecognitionClient: faceReconClient,
		ESClient:              esClient,
	}

	s.grpcPort = grpcPort
	s.db = db

	return s
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
	rows, err := s.db.Query(`SELECT id,  COALESCE(user_name, '') as user_name,
										COALESCE(real_name, '') as real_name,
										COALESCE(bio, '') as bio,
										COALESCE(avatar_url, '') as avatar_url
										FROM users WHERE LOWER(user_name) LIKE LOWER($1)`, fmt.Sprintf("%%%s%%", username.UserName))
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
	log.Println("writer is", s.userNamesWriter)
	if s.userNamesWriter != nil {
		log.Println("writing user", username.UserName, "to user topic")
		err := s.userNamesWriter.WriteMessages(context.Background(), kafka.Message{Value: []byte(username.UserName)})
		if err != nil {
			return nil, err
		}
	}

	u := &proto.User{}
	log.Println(username)

	err := s.db.QueryRow(`SELECT id, COALESCE(user_name, '') as user_name,
									COALESCE(real_name, '') as real_name,
									COALESCE(bio, '') as bio,
									COALESCE(avatar_url, '') as avatar_url
									FROM users WHERE user_name = $1`, username.UserName).Scan(&u.Id, &u.UserName, &u.RealName, &u.Bio, &u.AvatarUrl)
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

//GetUserWithUserId returns one User with the given username
func (s *GrpcServer) GetUserWithUserId(_ context.Context, username *proto.UserIdRequest) (*proto.User, error) {
	userID, err := strconv.ParseInt(username.UserId, 10, 64)
	if err != nil {
		return nil, err
	}

	u := &proto.User{}
	log.Println(username)

	err = s.db.QueryRow(`SELECT id, COALESCE(user_name, '') as user_name,
									COALESCE(real_name, '') as real_name,
									COALESCE(bio, '') as bio,
									COALESCE(avatar_url, '') as avatar_url
									FROM users WHERE id = $1`, userID).Scan(&u.Id, &u.UserName, &u.RealName, &u.Bio, &u.AvatarUrl)
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

	rows, err := s.db.Query(`SELECT id, COALESCE(post_id, '') as post_id,
										COALESCE(short_code, '') as short_code,
										COALESCE(caption, '') as caption,
										COALESCE(internal_picture_url, '') as internal_picture_url
										FROM posts WHERE user_id=$1`, request.UserId)
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

//GetTaggedPostsWithUserId returns all Posts the given User is tagged on
func (s *GrpcServer) GetTaggedPostsWithUserId(_ context.Context, request *proto.UserIdRequest) (*proto.InstaPostsResponse, error) {
	res := &proto.InstaPostsResponse{}

	rows, err := s.db.Query(`SELECT posts.id,
										COALESCE(posts.post_id, '') as post_id,
										COALESCE(posts.short_code, '') as short_code,
										COALESCE(posts.caption, '') as caption,
										COALESCE(posts.internal_picture_url, '') as internal_picture_url
										FROM posts
										JOIN post_tagged_users
										ON posts.id=post_tagged_users.post_id
										WHERE post_tagged_users.user_id=$1`, request.UserId)

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
	presignedURL, err := s.minioClient.PresignedGetObject(s.downloadBucketName, object, signedURLExpirationTime, reqParams)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return presignedURL.String(), nil

}

//SearchSimilarFaces to the given base64 encoded image
func (s *GrpcServer) SearchSimilarFaces(ctx context.Context, request *proto.FaceSearchRequest) (*proto.FaceSearchResponse, error) {
	imageContent, err := base64.StdEncoding.DecodeString(request.Base64EncodedPicture)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("couldn't base64 decode image: %w", err)
	}
	imagePath := utils.RandUUIDSeq()
	_, err = s.minioClient.PutObject(
		s.imageUploadBucket,
		imagePath,
		bytes.NewReader(imageContent),
		int64(len(imageContent)),
		minio.PutObjectOptions{},
	)

	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to upload image to S3: %w", err)
	}

	presignedURL, err := s.minioClient.PresignedGetObject(s.imageUploadBucket, imagePath, signedURLExpirationTime, url.Values{})
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("coudln't generate presgined URL: %w", err)
	}

	faces, err := s.facesClient.FindSimilarFacesInImage(presignedURL.String(), 10)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("coudln't find similar faces: %w", err)
	}

	if len(faces) == 0 {
		log.Println("no faces found")
		return &proto.FaceSearchResponse{Faces: nil}, nil
	}

	postIDs := []string{}

	for _, face := range faces {
		postIDs = append(postIDs, strconv.FormatInt(int64(face.PostID), 10))
	}
	sql := `SELECT id, COALESCE(internal_picture_url, '') as internal_picture_url FROM posts WHERE id in (` + strings.Join(postIDs, ",") + ")"
	log.Println(sql)
	rows, err := s.db.Query(sql)
	if err != nil {
		return nil, fmt.Errorf("failed to query postgres for posts: %w", err)
	}
	defer rows.Close()

	signedPostURLs := map[int]string{}
	for rows.Next() {
		var postID int
		var imgPath string
		rows.Scan(&postID, &imgPath)
		if imgPath == "" {
			continue
		}

		signedURL, err := s.getURLForPost(imgPath)
		if err != nil {
			return nil, fmt.Errorf("couldn't generate signed URL to post image: %w", err)
		}
		signedPostURLs[postID] = signedURL
	}

	responseFaces := []*proto.Face{}
	for _, foundFace := range faces {
		responseFaces = append(responseFaces, &proto.Face{
			PostId:       int32(foundFace.PostID),
			X:            int32(foundFace.X),
			Y:            int32(foundFace.Y),
			Width:        int32(foundFace.Width),
			Height:       int32(foundFace.Height),
			FullImageSrc: signedPostURLs[foundFace.PostID],
		})
	}

	return &proto.FaceSearchResponse{Faces: responseFaces}, nil
}

// SearchUsersWithWeightedPosts searches for users by their occurence in posts, taking the weights into account
func (s *GrpcServer) SearchUsersWithWeightedPosts(ctx context.Context, weightedPosts *proto.WeightedPosts) (*proto.WeightedUsers, error) {
	userIDToFaces := map[string]*proto.UserWithFaces{}
	for _, post := range weightedPosts.Posts {
		rows, err := s.db.Query(`SELECT u.id,  COALESCE(user_name, '') as user_name,
		COALESCE(u.real_name, '') as real_name,
		COALESCE(u.bio, '') as bio,
		COALESCE(u.avatar_url, '') as avatar_url
		FROM users u
		INNER JOIN posts p on p.user_id = u.id
		WHERE p.id = $1`, post.PostId)

		if err != nil {
			return nil, err
		}

		defer rows.Close()

		for rows.Next() {
			user := &proto.User{}

			err := rows.Scan(&user.Id, &user.UserName, &user.RealName, &user.Bio, &user.AvatarUrl)
			if err != nil {
				return nil, err
			}

			if userIDToFaces[user.Id] == nil {
				userIDToFaces[user.Id] = &proto.UserWithFaces{User: user}
			}

			wrappedUser := userIDToFaces[user.Id]

			wrappedUser.Faces = append(wrappedUser.Faces, post.Faces...)
			wrappedUser.Weight += post.Weight
		}
	}

	weightedUsers := &proto.WeightedUsers{}

	for _, user := range userIDToFaces {
		weightedUsers.UsersWithFaces = append(weightedUsers.UsersWithFaces, user)
	}

	return weightedUsers, nil
}

// DataPointCountForUserId counts all tables for the given user id
func (s *GrpcServer) DataPointCountForUserId(ctx context.Context, request *proto.UserIdRequest) (*proto.UserDataPointCount, error) {
	userID, err := strconv.ParseInt(request.UserId, 10, 64)
	if err != nil {
		return nil, err
	}

	row := s.db.QueryRow("select count(*) from posts where user_id = $1", userID)
	var postsCount int
	err = row.Scan(&postsCount)
	if err != nil {
		return nil, err
	}

	row = s.db.QueryRow("select count(*) from comments where owner_user_id = $1", userID)
	var commentsCount int
	err = row.Scan(&commentsCount)
	if err != nil {
		return nil, err
	}

	row = s.db.QueryRow("select count (*) from post_likes where user_id = $1", userID)
	var likesCount int
	err = row.Scan(&likesCount)
	if err != nil {
		return nil, err
	}

	row = s.db.QueryRow("select count(*) from follows where to_id = $1 or from_id = $1", userID)
	var followCount int
	err = row.Scan(&followCount)
	if err != nil {
		return nil, err
	}

	totalCount := postsCount + commentsCount + likesCount + followCount
	return &proto.UserDataPointCount{Count: int32(totalCount)}, nil
}

// FindCitiesForUserId in elasticsearch
func (s *GrpcServer) FindCitiesForUserId(ctx context.Context, request *proto.UserIdRequest) (*proto.FoundCities, error) {
	userID, err := strconv.ParseInt(request.UserId, 10, 64)
	if err != nil {
		return nil, err
	}

	foundCities := []string{}
	for city, cityTerms := range s.citiesMap {
		foundTerms, err := s.nlpAnalyzer.MatchTermsForUser(int(userID), cityTerms)
		if err != nil {
			panic(err)
		}
		log.Printf("city=%v \t-> foundTerms=%+v", city, foundTerms)
		// check if there are results for city
		if len(foundTerms) > 0 {
			foundCities = append(foundCities, city)
		}
	}
	log.Printf("foundCities=%+v", foundCities)

	return &proto.FoundCities{CityNames: foundCities}, nil
}
