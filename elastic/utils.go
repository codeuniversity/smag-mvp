package elastic

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
)

const (
	elasticProtocol = "http://%s"
)

// InitializeElasticSearch return setup elastic search client
func InitializeElasticSearch(esHosts []string) *elasticsearch.Client {

	var hosts []string
	for _, address := range esHosts {
		url := fmt.Sprintf(elasticProtocol, address)
		hosts = append(hosts, url)
	}

	cfg := elasticsearch.Config{
		Addresses: hosts,
	}
	client, err := elasticsearch.NewClient(cfg)

	if err != nil {
		panic(err)
	}
	return client
}
