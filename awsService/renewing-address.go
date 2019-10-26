package awsService

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	pb "github.com/codeuniversity/smag-mvp/awsService/proto"
	"github.com/codeuniversity/smag-mvp/models"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
	"log"
	"net"
)

// Scraper represents the scraper containing all clients it uses
type RenewingAddressGrpcServer struct {
	errQWriter *kafka.Writer
	awsSession *session.Session
	ec2Service *ec2.EC2
	grpcPort   int
	*service.Executor
}

// New returns an initilized scraper
func New(kafkaAddress string, grpcPort int) *RenewingAddressGrpcServer {
	var err error
	s := &RenewingAddressGrpcServer{}
	//s.errQWriter = kafka.NewWriter(kafka.WriterConfig{
	//	Brokers:  []string{kafkaAddress},
	//	Topic:    "renewing_elastic_ip_errors",
	//	Balancer: &kafka.LeastBytes{},
	//	Async:    false,
	//})
	s.awsSession, err = session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)
	if err != nil {
		fmt.Println("Error AWS Session ", err)
		return nil
	}
	s.ec2Service = ec2.New(s.awsSession)
	s.grpcPort = grpcPort
	return s
}

func (r *RenewingAddressGrpcServer) Listen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", r.grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	fmt.Println("Start RenewingAddressGrpcServer Server")

	grpcServer := grpc.NewServer()
	pb.RegisterRouteGuideServer(grpcServer, r)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (r *RenewingAddressGrpcServer) ErrorTest() {

	var err error
	err = &NoElasticIPError{"Test"}
	//err = errors.New("Test")

	var e *NoElasticIPError
	if !errors.As(err, &e) {
		fmt.Println("NoElasticIPError")
	} else {
		fmt.Println("Error")
	}
}

func (r *RenewingAddressGrpcServer) RenewElasticIp(context context.Context, reachedRequestLimit *pb.RenewingElasticIp) (*pb.RenewedElasticResult, error) {
	ec2Address, err := GetElasticPublicAddresses(r.ec2Service, reachedRequestLimit.InstanceId, reachedRequestLimit.PodIp)

	log.Println("PodIp: ", reachedRequestLimit.PodIp)
	log.Println("InstanceId: ", reachedRequestLimit.InstanceId)
	log.Println("Counter: ", reachedRequestLimit.Node)
	if err != nil {
		fmt.Println(err)
		return &pb.RenewedElasticResult{IsRenewed: false}, nil
	}

	var noElasticIPError *NoElasticIPError
	if !errors.As(err, &noElasticIPError) {
		err = disassociateAddress(r.awsSession, *ec2Address.PublicIp)
		if err != nil {
			fmt.Println("disassociateAddress Error: ", err)
			r.sendErrorMessage(reachedRequestLimit.InstanceId, err)
			return nil, err
		}

		err = releaseElasticAddresses(r.ec2Service, *ec2Address.AllocationId)
		if err != nil {
			fmt.Println("Error ReleaseElasticAddress: ", err)
			r.sendErrorMessage(reachedRequestLimit.InstanceId, err)
			return nil, err
		}
	}

	err = allocateAddresses(r.ec2Service, reachedRequestLimit.InstanceId, reachedRequestLimit.PodIp, *ec2Address.NetworkInterfaceId)
	if err != nil {
		fmt.Println("Error AllocateAddress: ", err)
		r.sendErrorMessage(reachedRequestLimit.InstanceId, err)
		return nil, err
	}

	return &pb.RenewedElasticResult{IsRenewed: true}, nil
}

func allocateAddresses(svc *ec2.EC2, instanceId string, localIp string, networkInterfaceId string) error {

	allocRes, err := svc.AllocateAddress(&ec2.AllocateAddressInput{
		Domain: aws.String("vpc"),
	})
	if err != nil {
		fmt.Println("AllocateAddress Error : ")
		return err
	}

	_, err = svc.AssociateAddress(&ec2.AssociateAddressInput{
		AllocationId:       allocRes.AllocationId,
		NetworkInterfaceId: &networkInterfaceId,
		PrivateIpAddress:   &localIp,
	})

	if err != nil {
		fmt.Println("AssociateAddress Error")
		return err
	}

	return nil
}

func releaseElasticAddresses(svc *ec2.EC2, allocationId string) error {

	_, err := svc.ReleaseAddress(&ec2.ReleaseAddressInput{
		AllocationId: &allocationId,
	})
	if err != nil {
		fmt.Println("Release Address Error ", err)
		return nil
	}

	fmt.Printf("Successfully released allocation ID\n")
	return nil
}

func GetElasticPublicAddresses(svc *ec2.EC2, instanceId string, localIp string) (*ec2.Address, error) {

	result, err := svc.DescribeAddresses(&ec2.DescribeAddressesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("instance-id"),
				Values: aws.StringSlice([]string{instanceId}),
			},
			{
				Name:   aws.String("private-ip-address"),
				Values: aws.StringSlice([]string{localIp}),
			},
		},
	})
	if err != nil {
		fmt.Println("Unable to elastic IP address, %v\n", err)
		return nil, err
	}

	if len(result.Addresses) == 0 {
		fmt.Printf("No elastic IPs for %s region\n", *svc.Config.Region)
		return nil, &NoElasticIPError{*svc.Config.Region}
	} else {
		fmt.Println("Elastic IPs")
		fmt.Println(len(result.Addresses))

		if len(result.Addresses) > 1 {
			//send Error Kafka Message!!!!
		}
		//publicIp := result.Addresses[0].PublicIp
		//networkInterfaceId := result.Addresses[0].NetworkInterfaceId
		return result.Addresses[0], nil
	}
}

// To disassociate an Elastic IP addresses in EC2-Classic
//
// This example disassociates an Elastic IP address from an instance in EC2-Classic.
func disassociateAddress(session *session.Session, publicIP string) error {

	svc := ec2.New(session)
	input := &ec2.DisassociateAddressInput{
		PublicIp: aws.String(publicIP),
	}

	_, err := svc.DisassociateAddress(input)
	if err != nil {
		return err
	}
	return nil
}

type NoElasticIPError struct {
	region string
}

func (e *NoElasticIPError) Error() string {
	return fmt.Sprintf("No elastic IPs for region %s \n", e.region)
}

func (r *RenewingAddressGrpcServer) sendErrorMessage(instanceId string, err error) {
	errMessage := &models.AwsServiceError{
		InstanceId: instanceId,
		Error:      err.Error(),
	}
	_, err = json.Marshal(errMessage)
	if err != nil {
		fmt.Println(err)
	}
	//r.errQWriter.WriteMessages(context.Background(), kafka.Message{Value: serializedErr})
}

func (r *RenewingAddressGrpcServer) Close() {
	r.errQWriter.Close()
}
