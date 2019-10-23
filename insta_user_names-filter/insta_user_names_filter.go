package filter

import (
	"encoding/json"

	"github.com/codeuniversity/smag-mvp/kafka/changestream"
	"github.com/segmentio/kafka-go"
)

type user struct {
	ID       int    `json:"id"`
	UserName string `json:"user_name"`
}

// UserNamesFilter represents the UserNamesFilter containing all clients it uses
type UserNamesFilter struct {
	*changestream.Filter
}

// New returns an initilized UserNamesFilter
func New(kafkaAddress, kafkaGroupID, changesTopic, userNameTopic string) *UserNamesFilter {
	t := &UserNamesFilter{
		Filter: changestream.NewFilter(kafkaAddress, kafkaGroupID, changesTopic, userNameTopic, filterChange),
	}

	return t
}

func filterChange(m *changestream.ChangeMessage) ([]kafka.Message, error) {
	if m.Payload.Op != "c" {
		return nil, nil
	}

	u := &user{}
	err := json.Unmarshal(m.Payload.After, u)
	if err != nil {
		return nil, err
	}

	return []kafka.Message{{Value: []byte(u.UserName)}}, nil
}
