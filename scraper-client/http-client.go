package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	pb "github.com/codeuniversity/smag-mvp/aws_service/proto"
	generator "github.com/codeuniversity/smag-mvp/http_header-generator"
	"github.com/codeuniversity/smag-mvp/utils"

	"io/ioutil"
	"net"
	"net/http"
	"time"

	"google.golang.org/grpc"
)

const (
	awsGetInstanceIdUrl = "http://169.254.169.254/latest/meta-data/instance-id"
)

type HttpClient struct {
	*generator.HTTPHeaderGenerator
	localAddressReachLimit bool
	localIp                string
	currentAddress         string
	client                 *http.Client
	instanceId             string
	grpcClient             *grpc.ClientConn
	scraperConfig          *ScraperConfig
}

func NewHttpClient(awsServiceAddress string, config *ScraperConfig) *HttpClient {
	client := &HttpClient{}
	client.HTTPHeaderGenerator = generator.New()
	client.scraperConfig = config
	var err error

	client.localIp = utils.MustGetStringFromEnv("POD_IP")
	client.client, err = client.getBoundAddressClient(client.localIp)

	if err != nil {
		panic(err)
	}

	client.instanceId, err = getAmazonInstanceId()

	if err != nil {
		log.Println("amazon InstanceId is not set")
	}

	client.grpcClient, err = grpc.Dial(awsServiceAddress, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	return client
}

func getAmazonInstanceId() (string, error) {
	resp, err := http.Get(awsGetInstanceIdUrl)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Println("Error: ", err)
		panic(err)
	}
	return string(body), err
}

func getPublicIp() string {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get("https://api.ipify.org?format=json")
	if err != nil {
		return ""
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Println("Error: ", err)
		return ""
	}

	type PublicIp struct {
		Ip string `json:"ip"`
	}
	var publicIp PublicIp
	err = json.Unmarshal(body, &publicIp)
	if err != nil {
		log.Println(err)
		return ""
	}
	return publicIp.Ip
}

func (h *HttpClient) getBoundAddressClient(localIp string) (*http.Client, error) {
	localAddr, err := net.ResolveIPAddr("ip", localIp)

	if err != nil {
		return nil, err
	}

	localTCPAddr := net.TCPAddr{
		IP: localAddr.IP,
	}

	d := net.Dialer{
		LocalAddr: &localTCPAddr,
		Timeout:   5 * time.Second,
		KeepAlive: 5 * time.Second,
	}

	tr := &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		DialContext:         d.DialContext,
		TLSHandshakeTimeout: 7 * time.Second,
	}

	return &http.Client{Transport: tr}, nil
}

func (h *HttpClient) WithRetries(requestRetryCount int, f func() error) error {
	var err error
	isWithRetriesDone := false
	var elasticIp string
	for i := 0; i < requestRetryCount; i++ {

		for k := 0; k < h.scraperConfig.ElasticIpRetryCount; k++ {
			err = f()
			ip := getPublicIp()
			log.Println("Real PublicIp: ", ip)
			log.Println("Elastic Ip: ", elasticIp)
			time.Sleep(time.Duration(h.scraperConfig.RequestTimeout) * time.Millisecond)

			if err == nil {
				return nil
			}
		}

		log.Println(err)
		elasticIp, err = h.checkIfIPReachedTheLimit(err)
		if err != nil {
			panic(err)
		}

		if !isWithRetriesDone && ((i + 1) == requestRetryCount) {
			isWithRetriesDone = true
			i--
		}
		time.Sleep(time.Duration(h.scraperConfig.ElasticAssignmentTimeout) * time.Millisecond)
	}
	return err
}

func (h *HttpClient) checkIfIPReachedTheLimit(err error) (string, error) {
	switch t := err.(type) {
	case *json.SyntaxError, *HTTPStatusError:
		log.Println("HttpStatus & Syntax Error ", t)
		elasticIp, err := h.sendRenewElasticIpRequestToAmazonService()
		if err != nil {
			return "", err
		}
		return elasticIp, nil
	default:
		return "", fmt.Errorf("error checkIfIPReachedTheLimit: %s", err)
	}
}

func (h *HttpClient) sendRenewElasticIpRequestToAmazonService() (string, error) {

	renewIp := pb.RenewingElasticIp{
		InstanceId: h.instanceId,
		PodIp:      h.localIp,
	}

	awsClient := pb.NewElasticIpServiceClient(h.grpcClient)
	result, err := awsClient.RenewElasticIp(context.Background(), &renewIp)
	if err != nil {
		log.Println("sendRenewElasticIpRequestToAmazonService Err: ", err)
		return "", err
	}

	return result.ElasticIp, nil
}

func (h *HttpClient) Close() {
	h.grpcClient.Close()
}

func (h *HttpClient) Do(request *http.Request) (*http.Response, error) {
	h.AddHeaders(&request.Header)
	return h.client.Do(request)
}
