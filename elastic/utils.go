package elastic

import (
	"github.com/elastic/go-elasticsearch/v7"
)

// InitializeElasticSearch returns an initialised elastic search client
func InitializeElasticSearch(esHosts []string) *elasticsearch.Client {
	cfg := elasticsearch.Config{
		Addresses: esHosts,
	}
	client, err := elasticsearch.NewClient(cfg)

	if err != nil {
		panic(err)
	}
	return client
}
