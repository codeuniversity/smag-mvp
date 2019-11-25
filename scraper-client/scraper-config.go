package client

import "github.com/codeuniversity/smag-mvp/utils"

type ScraperConfig struct {
	ElasticAssignmentTimeout int
	RequestTimeout           int
	RequestRetryCount        int
	ElasticIpRetryCount      int
}

func GetScraperConfig() *ScraperConfig {
	return &ScraperConfig{
		ElasticAssignmentTimeout: utils.GetNumberFromEnvWithDefault("ELASTIC_ASSIGNMENT_TIMEOUT", 10000),
		RequestTimeout:           utils.GetNumberFromEnvWithDefault("REQUEST_TIMEOUT", 1000),
		RequestRetryCount:        utils.GetNumberFromEnvWithDefault("REQUEST_RETRY_COUNT", 3),
		ElasticIpRetryCount:      utils.GetNumberFromEnvWithDefault("ELASTIC_IP_RETRY_COUNT", 2),
	}
}
