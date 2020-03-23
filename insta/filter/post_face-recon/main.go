package main

import (
	"encoding/json"

	"github.com/codeuniversity/smag-mvp/insta/models"
	"github.com/codeuniversity/smag-mvp/kafka/changestream"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/codeuniversity/smag-mvp/utils"

	"github.com/segmentio/kafka-go"
)

func main() {
	kafkaAddress := utils.GetStringFromEnvWithDefault("KAFKA_ADDRESS", "127.0.0.1:9092")
	groupID := utils.MustGetStringFromEnv("KAFKA_GROUPID")
	changesTopic := utils.GetStringFromEnvWithDefault("KAFKA_CHANGE_TOPIC", "postgres.public.posts")
	downloadTopic := utils.GetStringFromEnvWithDefault("KAFKA_PICTURE_FACE_RECON_TOPIC", "insta_post_face_recon_jobs")

	f := changestream.NewFilter(kafkaAddress, groupID, changesTopic, downloadTopic, filterChange)

	service.CloseOnSignal(f)
	waitUntilClosed := f.Start()

	waitUntilClosed()
}

type post struct {
	ID                 int    `json:"id"`
	InternalPictureURL string `json:"internal_picture_url"`
}

func filterChange(m *changestream.ChangeMessage) ([]kafka.Message, error) {
	if !(m.Payload.Op == "c" || m.Payload.Op == "u") {
		return nil, nil
	}

	currentVersion := &post{}
	err := json.Unmarshal(m.Payload.After, currentVersion)
	if err != nil {
		return nil, err
	}

	if m.Payload.Op == "c" {
		return constructDownloadJobMessage(currentVersion)
	}

	previousVersion := &post{}
	err = json.Unmarshal(m.Payload.Before, previousVersion)
	if err != nil {
		return nil, err
	}

	if currentVersion.InternalPictureURL != previousVersion.InternalPictureURL {
		return constructDownloadJobMessage(currentVersion)
	}

	return nil, nil
}

func constructDownloadJobMessage(p *post) ([]kafka.Message, error) {
	if p.InternalPictureURL == "" {
		return nil, nil
	}

	job := &models.PostDownloadJob{
		PostID:     p.ID,
		PictureURL: p.InternalPictureURL,
	}
	b, err := json.Marshal(job)
	if err != nil {
		return nil, err
	}

	return []kafka.Message{
		{Value: b},
	}, nil
}
