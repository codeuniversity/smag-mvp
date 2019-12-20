package aws_service

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	pb "github.com/codeuniversity/smag-mvp/aws_service/proto"
	"google.golang.org/grpc"
)

// Scraper represents the scraper containing all clients it uses
type RenewingAddressGrpcServer struct {
	awsSession *session.Session
	ec2Service *ec2.EC2
	grpcPort   string
}

// New returns an initialised gRPC Server
func New(grpcPort string) *RenewingAddressGrpcServer {
	var err error
	s := &RenewingAddressGrpcServer{}
	s.awsSession, err = session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)
	if err != nil {
		log.Println("Error AWS Session ", err)
		return nil
	}
	s.ec2Service = ec2.New(s.awsSession)
	s.grpcPort = grpcPort
	return s
}

// GrpcServer represents the gRPC Server containing the db connection and port
func (r *RenewingAddressGrpcServer) Listen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", r.grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterElasticIpServiceServer(grpcServer, r)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// renewing the elastic ip using local ip and instanceId
func (r *RenewingAddressGrpcServer) RenewElasticIp(context context.Context, reachedRequestLimit *pb.RenewingElasticIp) (*pb.RenewedElasticResult, error) {

	log.Println("PodIp: ", reachedRequestLimit.PodIp)
	log.Println("InstanceId: ", reachedRequestLimit.InstanceId)
	var networkInterfaceId string
	ec2Address, err := r.getElasticPublicAddresses(reachedRequestLimit.InstanceId, reachedRequestLimit.PodIp)

	if err != nil {
		log.Println(err)
		networkInterfaceId, err = r.getNetworkInterfaceId(reachedRequestLimit.InstanceId, reachedRequestLimit.PodIp)

		if err != nil {
			return nil, err
		}
	} else {
		networkInterfaceId = *ec2Address.NetworkInterfaceId
		err = r.disassociateAddress(*ec2Address.PublicIp)
		if err != nil {
			log.Println("Error DisassociateAddress: ", err)
			return nil, err
		}

		err = r.releaseElasticAddresses(*ec2Address.AllocationId)
		if err != nil {
			log.Println("Error ReleaseElasticAddress: ", err)
			return nil, err
		}
	}

	elasticIp, err := allocateAddresses(r.ec2Service, reachedRequestLimit.PodIp, networkInterfaceId)
	if err != nil {
		log.Println("Error AllocateAddress: ", err)
		return nil, err
	}

	return &pb.RenewedElasticResult{ElasticIp: elasticIp}, nil
}

func (r *RenewingAddressGrpcServer) getNetworkInterfaceId(instanceId string, localIp string) (string, error) {
	networkInterfaces, err := r.ec2Service.DescribeNetworkInterfaces(&ec2.DescribeNetworkInterfacesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("attachment.instance-id"),
				Values: aws.StringSlice([]string{instanceId}),
			},
			{
				Name:   aws.String("addresses.private-ip-address"),
				Values: aws.StringSlice([]string{localIp}),
			},
		},
	})

	if err != nil {
		return "", err
	}

	if len(networkInterfaces.NetworkInterfaces) == 0 {
		return "", fmt.Errorf("No network interface Id %s", *r.ec2Service.Config.Region)
	}
	return *networkInterfaces.NetworkInterfaces[0].NetworkInterfaceId, nil
}

func allocateAddresses(svc *ec2.EC2, localIp string, networkInterfaceId string) (string, error) {

	allocRes, err := svc.AllocateAddress(&ec2.AllocateAddressInput{
		Domain: aws.String("vpc"),
	})

	if err != nil {
		log.Println("AllocateAddress Error : ")
		return "", err
	}

	_, err = svc.AssociateAddress(&ec2.AssociateAddressInput{
		AllocationId:       allocRes.AllocationId,
		NetworkInterfaceId: &networkInterfaceId,
		PrivateIpAddress:   &localIp,
	})

	if err != nil {
		log.Println("AssociateAddress Error")
		return "", err
	}

	return *allocRes.PublicIp, nil
}

func (r *RenewingAddressGrpcServer) releaseElasticAddresses(allocationId string) error {

	_, err := r.ec2Service.ReleaseAddress(&ec2.ReleaseAddressInput{
		AllocationId: &allocationId,
	})

	if err != nil {
		log.Println("Release Address Error ", err)
		return err
	}

	log.Printf("Successfully released allocation ID\n")
	return nil
}

func (r *RenewingAddressGrpcServer) getElasticPublicAddresses(instanceId string, localIp string) (*ec2.Address, error) {

	result, err := r.ec2Service.DescribeAddresses(&ec2.DescribeAddressesInput{
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
		log.Println("Unable to get elastic IP address: ", err)
		return nil, err
	}

	if len(result.Addresses) == 0 {
		log.Printf("No elastic IPs for %s region\n", *r.ec2Service.Config.Region)
		return nil, &NoElasticIPError{*r.ec2Service.Config.Region}
	}
	return result.Addresses[0], nil
}

// To disassociate an Elastic IP addresses in EC2-Classic
//
// This example disassociates an Elastic IP address from an instance in EC2-Classic.
func (r *RenewingAddressGrpcServer) disassociateAddress(publicIP string) error {
	svc := ec2.New(r.awsSession)
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
