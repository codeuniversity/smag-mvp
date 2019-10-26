package scraper_client

import (
	"context"
	"encoding/json"
	"fmt"
	pb "github.com/codeuniversity/smag-mvp/aws_service/proto"
	generator "github.com/codeuniversity/smag-mvp/http_header-generator"
	"github.com/codeuniversity/smag-mvp/utils"

	"google.golang.org/grpc"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

type HttpClient struct {
	*generator.HTTPHeaderGenerator
	localAddressReachLimit bool
	localIp                string
	currentAddress         string
	client                 *http.Client
	instanceId             string
	grpcClient             *grpc.ClientConn
}

func NewHttpClient(awsServiceAddress string) *HttpClient {
	client := &HttpClient{}
	client.HTTPHeaderGenerator = generator.New()
	var err error

	client.localIp = utils.MustGetStringFromEnv("POD_IP")
	client.client, err = client.getBoundAddressClient(client.localIp)

	if err != nil {
		panic(err)
	}

	client.instanceId, err = getAmazonInstanceId()

	if err != nil {
		fmt.Println("amazon InstanceId is not set")
	}

	client.grpcClient, err = grpc.Dial(awsServiceAddress, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	return client
}

func getAmazonInstanceId() (string, error) {
	resp, err := http.Get("http://169.254.169.254/latest/meta-data/instance-id")
	if err != nil {
		return "", nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("Error: ", err)
		panic(err)
	}
	return string(body), nil
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
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	tr := &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		DialContext:         d.DialContext,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	return &http.Client{Transport: tr}, nil
}

func (h *HttpClient) WithRetries(times int, f func() error) error {
	var err error
	for i := 0; i < times; i++ {
		err = f()
		if err == nil {
			return nil
		}

		fmt.Println(err)
		isRenewed, err := h.checkIfIPReachedTheLimit(err)
		if err != nil {
			fmt.Println(err)
		}
		if isRenewed {
			times++
		}
		time.Sleep(380 * time.Millisecond)
	}
	return err
}

func (h *HttpClient) checkIfIPReachedTheLimit(err error) (bool, error) {
	fmt.Println("checkIfIPReachedTheLimit")
	switch t := err.(type) {
	case *json.SyntaxError, *HTTPStatusError:
		fmt.Println("SyntaxError")
		_, err := h.sendRenewElasticIpRequestToAmazonService()
		if err != nil {
			return false, err
		}
		return true, nil
	default:
		fmt.Println("Found Wrong Json Type Error ", t)
		return false, err
	}
}

func (h *HttpClient) sendRenewElasticIpRequestToAmazonService() (bool, error) {

	renewIp := pb.RenewingElasticIp{
		InstanceId: h.instanceId,
		Node:       "",
		Pod:        "",
		PodIp:      h.localIp,
	}

	awsClient := pb.NewRouteGuideClient(h.grpcClient)
	result, err := awsClient.RenewElasticIp(context.Background(), &renewIp)
	if err != nil {
		fmt.Println("sendRenewElasticIpRequestToAmazonService Err: ", err)
		return false, err
	}

	return result.IsRenewed, nil
}

func (h *HttpClient) Close() {
	h.grpcClient.Close()
}

func (h *HttpClient) Do(request *http.Request) (*http.Response, error) {
	h.AddHeaders(&request.Header)
	return h.client.Do(request)
}
