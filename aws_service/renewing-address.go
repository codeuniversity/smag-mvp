package aws_service

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	pb "github.com/codeuniversity/smag-mvp/aws_service/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

// Scraper represents the scraper containing all clients it uses
type RenewingAddressGrpcServer struct {
	awsSession *session.Session
	ec2Service *ec2.EC2
	grpcPort   int
}

// New returns an initialised gRPC Server
func New(grpcPort int) *RenewingAddressGrpcServer {
	var err error
	s := &RenewingAddressGrpcServer{}
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
	pb.RegisterElasticIpServiceServer(grpcServer, r)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (r *RenewingAddressGrpcServer) RenewElasticIp(context context.Context, reachedRequestLimit *pb.RenewingElasticIp) (*pb.RenewedElasticResult, error) {
	ec2Address, err := r.getElasticPublicAddresses(reachedRequestLimit.InstanceId, reachedRequestLimit.PodIp)

	log.Println("PodIp: ", reachedRequestLimit.PodIp)
	log.Println("InstanceId: ", reachedRequestLimit.InstanceId)
	if err != nil {
		fmt.Println(err)
		return &pb.RenewedElasticResult{IsRenewed: false}, err
	}

	err = r.disassociateAddress(*ec2Address.PublicIp)
	if err != nil {
		fmt.Println("disassociateAddress Error: ", err)
		return nil, err
	}

	err = r.releaseElasticAddresses(*ec2Address.AllocationId)
	if err != nil {
		fmt.Println("Error ReleaseElasticAddress: ", err)
		return nil, err
	}

	err = allocateAddresses(r.ec2Service, reachedRequestLimit.PodIp, *ec2Address.NetworkInterfaceId)
	if err != nil {
		fmt.Println("Error AllocateAddress: ", err)
		return nil, err
	}

	return &pb.RenewedElasticResult{IsRenewed: true}, nil
}

func allocateAddresses(svc *ec2.EC2, localIp string, networkInterfaceId string) error {

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

func (r *RenewingAddressGrpcServer) releaseElasticAddresses(allocationId string) error {

	_, err := r.ec2Service.ReleaseAddress(&ec2.ReleaseAddressInput{
		AllocationId: &allocationId,
	})
	if err != nil {
		fmt.Println("Release Address Error ", err)
		return err
	}

	fmt.Printf("Successfully released allocation ID\n")
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
		fmt.Println("Unable to elastic IP address: ", err)
		return nil, err
	}

	if len(result.Addresses) == 0 {
		fmt.Printf("No elastic IPs for %s region\n", *r.ec2Service.Config.Region)
		return nil, &NoElasticIPError{*r.ec2Service.Config.Region}
	} else {
		return result.Addresses[0], nil
	}
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
