package awsService

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/codeuniversity/smag-mvp/models"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/segmentio/kafka-go"
	"time"
)

// Scraper represents the scraper containing all clients it uses
type RenewingAddresses struct {
	reachedLimitQReader   *kafka.Reader
	renewedAddressQWriter *kafka.Writer
	errQWriter            *kafka.Writer
	awsSession            *session.Session
	ec2Service            *ec2.EC2
	*service.Executor
}

// New returns an initilized scraper
func New(kafkaAddress string) *RenewingAddresses {
	var err error
	s := &RenewingAddresses{}
	s.reachedLimitQReader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{kafkaAddress},
		GroupID:        "renewing_addresses_group",
		Topic:          "reached_limit",
		CommitInterval: time.Minute * 10,
	})
	s.renewedAddressQWriter = kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{kafkaAddress},
		Topic:    "renewed_elastic_ip",
		Balancer: &kafka.LeastBytes{},
		Async:    true,
	})
	s.errQWriter = kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{kafkaAddress},
		Topic:    "renewing_elastic_ip_errors",
		Balancer: &kafka.LeastBytes{},
		Async:    false,
	})
	s.Executor = service.New()
	s.awsSession, err = session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)
	if err != nil {
		fmt.Println("Error AWS Session ", err)
		return nil
	}
	s.ec2Service = ec2.New(s.awsSession)
	return s
}

func (r *RenewingAddresses) Run() {
	defer func() {
		r.MarkAsStopped()
	}()

	fmt.Println("starting AWS-Service")
	for r.IsRunning() {
		m, err := r.reachedLimitQReader.FetchMessage(context.Background())
		fmt.Println("Message Time: ", m.Time)
		fmt.Println("START")
		if err != nil {
			fmt.Println(err)
			break
		}

		var reachedRequestLimit models.RenewingAddresses
		err = json.Unmarshal(m.Value, &reachedRequestLimit)
		if err != nil {
			fmt.Println("Json: ", err)
		}

		publicIP, allocationIds, err := getElasticPublicAddresses(r.ec2Service, reachedRequestLimit.InstanceId)

		if err != nil {
			fmt.Println("getElasticPublicAddress Error: ", err)
			r.sendErrorMessage(m, reachedRequestLimit.InstanceId, err)
			return
		}
		err = disassociateAddress(r.awsSession, publicIP)
		if err != nil {
			fmt.Println("disassociateAddress Error: ", err)
			r.sendErrorMessage(m, reachedRequestLimit.InstanceId, err)
			return
		}

		err = releaseElasticAddresses(r.ec2Service, allocationIds)
		if err != nil {
			fmt.Println("Error ReleaseElasticAddress: ", err)
			r.sendErrorMessage(m, reachedRequestLimit.InstanceId, err)
			return
		}

		err = allocateAddresses(r.ec2Service, reachedRequestLimit.InstanceId, 2, reachedRequestLimit.LocalIps)
		if err != nil {
			fmt.Println("Error AllocateAddress: ", err)
			r.sendErrorMessage(m, reachedRequestLimit.InstanceId, err)
			return
		}

		err = r.renewedAddressQWriter.WriteMessages(context.Background(), kafka.Message{Value: m.Value})

		if err != nil {
			fmt.Println(err)
			break
		}

		r.reachedLimitQReader.CommitMessages(context.Background(), m)
	}
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

func (r *RenewingAddresses) sendErrorMessage(m kafka.Message, instanceId string, err error) {
	errMessage := &models.AwsServiceError{
		InstanceId: instanceId,
		Error:      err.Error(),
	}
	serializedErr, err := json.Marshal(errMessage)
	if err != nil {
		fmt.Println(err)
	}
	r.errQWriter.WriteMessages(context.Background(), kafka.Message{Value: serializedErr})
	r.reachedLimitQReader.CommitMessages(context.Background(), m)
}

func (r *RenewingAddresses) Close() {
	r.Stop()
	r.WaitUntilStopped(time.Second * 3)

	r.errQWriter.Close()
	r.reachedLimitQReader.Close()
	r.renewedAddressQWriter.Close()
	r.MarkAsClosed()
}
