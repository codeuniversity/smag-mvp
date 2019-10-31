package main

import (
	"testing"

	"github.com/codeuniversity/smag-mvp/kafka/changestream"
	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	validPayloadJSON := []byte("{\"id\":1,\"picture_url\":\"https://test.test\"}")
	preUpdatePayloadJSON := []byte("{\"id\":1,\"picture_url\":\"https://test2.test2\"}")
	invalidPayloadJSON := []byte("{\"id\":\"1\",\"picture_url\":\"https://test.test\"}")

	t.Run("create event with unmarshable json", func(t *testing.T) {
		//create test input
		changeMsg := &changestream.ChangeMessage{}
		changeMsg.Payload.Op = "c"
		changeMsg.Payload.After = validPayloadJSON

		kMessages, err := filterChange(changeMsg)

		expected := "{\"post_id\":1,\"picture_url\":\"https://test.test\"}"

		assert.Nil(t, err, "no error")
		assert.Equal(t, 1, len(kMessages))
		assert.Equal(t, expected, string(kMessages[0].Value))
	})

	t.Run("update event with unmarshable json", func(t *testing.T) {
		//create test input
		changeMsg := &changestream.ChangeMessage{}
		changeMsg.Payload.Op = "u"
		changeMsg.Payload.Before = preUpdatePayloadJSON
		changeMsg.Payload.After = validPayloadJSON

		kMessages, err := filterChange(changeMsg)

		expected := "{\"post_id\":1,\"picture_url\":\"https://test.test\"}"

		assert.Nil(t, err, "no error")
		assert.Equal(t, 1, len(kMessages))
		assert.Equal(t, expected, string(kMessages[0].Value))
	})

	t.Run("create event with not unmarshable json", func(t *testing.T) {
		//create test input
		changeMsg := &changestream.ChangeMessage{}
		changeMsg.Payload.Op = "c"
		changeMsg.Payload.After = invalidPayloadJSON

		kMessages, err := filterChange(changeMsg)

		assert.NotNil(t, err, "error occurs")
		assert.Nil(t, kMessages, "nil output")
	})

	t.Run("ignored event", func(t *testing.T) {
		//create test input
		changeMsg := &changestream.ChangeMessage{}
		changeMsg.Payload.Op = "d"

		kMessages, err := filterChange(changeMsg)

		assert.Nil(t, err, "no error")
		assert.Nil(t, kMessages, "nil output")
	})
}
