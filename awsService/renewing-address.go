package awsService

import (
	"context"
	"encoding/json"
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
	s.errQWriter = kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{kafkaAddress},
		Topic:    "renewing_elastic_ip_errors",
		Balancer: &kafka.LeastBytes{},
		Async:    false,
	})
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

func (r *RenewingAddressGrpcServer) RenewElasticIp(context context.Context, reachedRequestLimit *pb.RenewingElasticIp) (*pb.RenewedElasticResult, error) {
	publicIP, allocationIds, err := getElasticPublicAddresses(r.ec2Service, reachedRequestLimit.InstanceId)

	if err != nil {
		fmt.Println("getElasticPublicAddress Error: ", err)
		r.sendErrorMessage(reachedRequestLimit.InstanceId, err)
		return nil, err
	}
	err = disassociateAddress(r.awsSession, publicIP)
	if err != nil {
		fmt.Println("disassociateAddress Error: ", err)
		r.sendErrorMessage(reachedRequestLimit.InstanceId, err)
		return nil, err
	}

	err = releaseElasticAddresses(r.ec2Service, allocationIds)
	if err != nil {
		fmt.Println("Error ReleaseElasticAddress: ", err)
		r.sendErrorMessage(reachedRequestLimit.InstanceId, err)
		return nil, err
	}

	localIps := []string{reachedRequestLimit.PodIp}
	err = allocateAddresses(r.ec2Service, reachedRequestLimit.InstanceId, 2, localIps)
	if err != nil {
		fmt.Println("Error AllocateAddress: ", err)
		r.sendErrorMessage(reachedRequestLimit.InstanceId, err)
		return nil, err
	}

	return &pb.RenewedElasticResult{IsRenewed: true}, nil
}

func allocateAddresses(svc *ec2.EC2, instanceId string, count int, localIps []string) error {

	result, err := svc.DescribeNetworkInterfaces(&ec2.DescribeNetworkInterfacesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("attachment.instance-id"),
				Values: aws.StringSlice([]string{instanceId}),
			},
		},
	})

	if err != nil {
		return err
	}

	foundLocalIpCounter := 0
	for _, networkInterface := range result.NetworkInterfaces {
		for _, ip := range networkInterface.PrivateIpAddresses {
			for _, localAddress := range localIps {
				if localAddress == *ip.PrivateIpAddress {
					foundLocalIpCounter++

					allocRes, err := svc.AllocateAddress(&ec2.AllocateAddressInput{
						Domain: aws.String("vpc"),
					})
					if err != nil {
						fmt.Println("AllocateAddress Error : ")
						return err
					}

					_, err = svc.AssociateAddress(&ec2.AssociateAddressInput{
						AllocationId:       allocRes.AllocationId,
						NetworkInterfaceId: networkInterface.NetworkInterfaceId,
						PrivateIpAddress:   &localAddress,
					})

					if err != nil {
						fmt.Println("AssociateAddress Error")
						return err
					}
				}
			}
		}
	}
	fmt.Println("FoundLocalIpCounter: ", foundLocalIpCounter)
	return nil
}

func releaseElasticAddresses(svc *ec2.EC2, allocationIds []string) error {

	for _, allocationId := range allocationIds {
		_, err := svc.ReleaseAddress(&ec2.ReleaseAddressInput{
			AllocationId: &allocationId,
		})
		if err != nil {
			fmt.Println("Release Address Error ", err)
			return nil
		}

	}

	fmt.Printf("Successfully released allocation ID\n")
	return nil
}

func getElasticPublicAddresses(svc *ec2.EC2, instanceId string) ([]string, []string, error) {

	result, err := svc.DescribeAddresses(&ec2.DescribeAddressesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("instance-id"),
				Values: aws.StringSlice([]string{instanceId}),
			},
		},
	})
	if err != nil {
		fmt.Println("Unable to elastic IP address, %v", err)
		return nil, nil, err
	}

	if len(result.Addresses) == 0 {
		fmt.Printf("No elastic IPs for %s region\n", *svc.Config.Region)
		return nil, nil, &NoElasticIPError{*svc.Config.Region}
	} else {
		fmt.Println("Elastic IPs")
		var publicIp []string
		var allocationIds []string
		for _, addr := range result.Addresses {
			allocationIds = append(allocationIds, *addr.AllocationId)
			publicIp = append(publicIp, aws.StringValue(addr.PublicIp))
		}
		return publicIp, allocationIds, nil
	}
}

// To disassociate an Elastic IP addresses in EC2-Classic
//
// This example disassociates an Elastic IP address from an instance in EC2-Classic.
func disassociateAddress(session *session.Session, publicIPs []string) error {

	for _, ip := range publicIPs {
		svc := ec2.New(session)
		input := &ec2.DisassociateAddressInput{
			PublicIp: aws.String(ip),
		}

		_, err := svc.DisassociateAddress(input)
		if err != nil {
			return err
		}
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
	serializedErr, err := json.Marshal(errMessage)
	if err != nil {
		fmt.Println(err)
	}
	r.errQWriter.WriteMessages(context.Background(), kafka.Message{Value: serializedErr})
}

func (r *RenewingAddressGrpcServer) Close() {
	r.errQWriter.Close()
}
