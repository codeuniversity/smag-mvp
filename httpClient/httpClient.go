package httpClient

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/codeuniversity/smag-mvp/models"
	"github.com/segmentio/kafka-go"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"time"
)

var userAccountInfoUrl = "https://instagram.com/%s/?__a=1"
var userAccountMediaUrl = "https://www.instagram.com/graphql/query/?query_hash=58b6785bea111c67129decbe6a448951&variables=%s"

type HttpClient struct {
	browserAgent             BrowserAgent
	localAddressesReachLimit map[string]bool
	currentAddress           string
	client                   *http.Client
	renewedAddressQReader    *kafka.Reader
	reachedLimitQWriter      *kafka.Writer
	instanceId               string
}

func New(localAddressCount int, kafkaAddress string) *HttpClient {
	client := &HttpClient{}
	client.renewedAddressQReader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{kafkaAddress},
		GroupID:        "renewed_elastic_ip",
		Topic:          "user_names",
		CommitInterval: time.Minute * 10,
	})

	client.reachedLimitQWriter = kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{kafkaAddress},
		Topic:    "reached_limit",
		Balancer: &kafka.LeastBytes{},
		Async:    true,
	})

	data, err := ioutil.ReadFile("useragents.json")
	if err != nil {
		panic(err)
	}
	var userAgent BrowserAgent
	errJson := json.Unmarshal(data, &userAgent)

	if errJson != nil {
		panic(errJson)
	}
	client.browserAgent = userAgent
	client.localAddressesReachLimit = make(map[string]bool)
	addresses, err := getLocalIpAddresses(localAddressCount)
	if err != nil {
		fmt.Println("httpClient getLocalIpAddresses Error: ", err)
		return nil
	}

	for _, localIp := range addresses {
		client.localAddressesReachLimit[localIp] = true
	}

	client.client, err = client.getClient(addresses[0])

	if err != nil {
		panic(err)
	}

	client.instanceId = getAmazonInstanceId()
	return client
}

func getAmazonInstanceId() string {
	resp, err := http.Get("http://169.254.169.254/latest/meta-data/instance-id")
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("Error: ", err)
		panic(err)
	}
	return string(body)
}

func getLocalIpAddresses(count int) ([]string, error) {
	interfaces, err := net.Interfaces()

	if err != nil {
		panic(err)
	}

	var localAddresses []string
	for _, networkInterface := range interfaces {
		fmt.Println("NetworkInterface: ", networkInterface.Name)
		addrs, err := networkInterface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, address := range addrs {
			localAddresses = append(localAddresses, address.String())
		}
	}

	if len(localAddresses) >= count {
		panic(fmt.Sprintf("Not Enough Local Ip Addresses, Requirement: %d \n", count))
	}

	return localAddresses[:count], nil
}

func (h *HttpClient) getClient(localIp string) (*http.Client, error) {
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

func (h *HttpClient) ScrapeAccountInfo(username string) (models.InstagramAccountInfo, error) {
	var userAccountInfo models.InstagramAccountInfo
	url := fmt.Sprintf(userAccountInfoUrl, username)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return userAccountInfo, err
	}
	h.getHeaders(request)

	response, err := h.client.Do(request)
	if err != nil {
		return userAccountInfo, err
	}
	if response.StatusCode != 200 {
		return userAccountInfo, &HttpStatusError{fmt.Sprintf("Error HttpStatus: %s", response.StatusCode)}
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return userAccountInfo, err
	}
	err = json.Unmarshal(body, &userAccountInfo)
	if err != nil {
		return userAccountInfo, err
	}
	return userAccountInfo, nil
}

func (h *HttpClient) ScrapeProfileMedia(userId string, endCursor string) (models.InstagramMedia, error) {
	var instagramMedia models.InstagramMedia

	type Variables struct {
		Id    string `json:"id"`
		First int    `json:"first"`
		After string `json:"after"`
	}
	variable := &Variables{userId, 12, endCursor}
	variableJson, err := json.Marshal(variable)
	fmt.Println(string(variableJson))
	if err != nil {
		return instagramMedia, err
	}
	queryEncoded := url.QueryEscape(string(variableJson))
	url := fmt.Sprintf(userAccountMediaUrl, queryEncoded)

	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return instagramMedia, err
	}
	response, err := h.client.Do(request)
	if err != nil {
		return instagramMedia, err
	}
	if response.StatusCode != 200 {
		return instagramMedia, &HttpStatusError{fmt.Sprintf("Error HttpStatus: %s", response.StatusCode)}
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return instagramMedia, err
	}
	err = json.Unmarshal(body, &instagramMedia)
	if err != nil {
		return instagramMedia, err
	}
	return instagramMedia, nil
}

func (h *HttpClient) WithRetries(times int, f func() error) error {
	var err error
	for i := 0; i < times; i++ {
		err = f()
		if err == nil {
			return nil
		}

		fmt.Println(err)
		foundAddress, err := h.checkIfIPReachedTheLimit(err)

		if err != nil {
			fmt.Println(err)
		}
		if foundAddress {
			times++
		}
		time.Sleep(400 * time.Millisecond)
	}
	return err
}

func (h *HttpClient) checkIfIPReachedTheLimit(err error) (bool, error) {
	switch t := err.(type) {
	case *json.SyntaxError:
		addresses, foundAddress := h.checkAvailableAddresses()

		if foundAddress {
			return true, nil
		}
		if h.localAddressesReachLimit[h.currentAddress] == false {
			err := h.sendRenewElasticIpRequestToAmazonService(addresses)
			if err != nil {
				return false, err
			}

			renewedAddresses := models.RenewingAddresses{InstanceId: "", LocalIps: []string{}}

			for renewedAddresses.InstanceId != h.instanceId {

				renewedAddresses, err := h.waitForRenewElasticIpRequest()
				if err != nil {
					return false, err
				}

				if renewedAddresses.InstanceId == h.instanceId {
					for ip := range h.localAddressesReachLimit {
						h.localAddressesReachLimit[ip] = true
					}
					return true, nil
				}
			}
		}
	default:
		fmt.Println("Found Wrong Json Type Error ", t)
		return false, err
	}
	fmt.Println("checkIfIPReachedTheLimit is not working!!!")
	return false, err
}

func (h *HttpClient) checkAvailableAddresses() ([]string, bool) {
	h.localAddressesReachLimit[h.currentAddress] = false
	var addresses []string
	var err error
	for ip := range h.localAddressesReachLimit {
		addresses = append(addresses, ip)
		if h.localAddressesReachLimit[ip] == true {
			h.currentAddress = ip
			h.client, err = h.getClient(ip)
			if err != nil {
				panic(err)
			}
			return addresses, true
		}
	}
	return addresses, false
}
func (h *HttpClient) sendRenewElasticIpRequestToAmazonService(addresses []string) error {
	renewAddresses := models.RenewingAddresses{
		InstanceId: h.instanceId,
		LocalIps:   addresses,
	}

	renewAdressesJson, err := json.Marshal(renewAddresses)
	if err != nil {
		return err
	}
	h.reachedLimitQWriter.WriteMessages(context.Background(), kafka.Message{Value: renewAdressesJson})
	return nil
}

func (h *HttpClient) waitForRenewElasticIpRequest() (*models.RenewingAddresses, error) {
	message, err := h.renewedAddressQReader.FetchMessage(context.Background())

	if err != nil {
		return nil, err
	}

	var renewedAddresses models.RenewingAddresses
	err = json.Unmarshal(message.Value, &renewedAddresses)
	if err != nil {
		return nil, err
	}

	return &renewedAddresses, err
}

type HttpStatusError struct {
	s string
}

func (e *HttpStatusError) Error() string {
	return e.s
}

type BrowserAgent []struct {
	UserAgents string `json:"useragent"`
}

func (h *HttpClient) getRandomUserAgent() string {
	randomNumber := rand.Intn(len(h.browserAgent))
	return h.browserAgent[randomNumber].UserAgents
}

func (h *HttpClient) Close() {
	h.renewedAddressQReader.Close()
	h.reachedLimitQWriter.Close()
}

func (h *HttpClient) getHeaders(request *http.Request) {
	request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3")
	request.Header.Add("Accept-Charset", "utf-8")
	request.Header.Add("Accept-Language", "de-DE,de;q=0.9,en-US;q=0.8,en;q=0.7")
	request.Header.Add("Cache-Control", "no-cache")
	request.Header.Add("Content-Type", "application/json; charset=utf-8")
	request.Header.Add("User-Agent", h.getRandomUserAgent())
}
