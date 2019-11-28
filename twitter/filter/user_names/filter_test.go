package main

import (
	"testing"

	"github.com/codeuniversity/smag-mvp/kafka/changestream"
	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	validPayloadJSON := []byte("{\"id\":1,\"username\":\"TestUser\"}")
	invalidPayloadJSON := []byte("{\"id\":\"1\",\"username\":\"TestUser\"}")

	t.Run("create event with unmarshable json", func(t *testing.T) {
		//create test input
		changeMsg := &changestream.ChangeMessage{}
		changeMsg.Payload.Op = "c"
		changeMsg.Payload.After = validPayloadJSON

		kMessages, err := filterChange(changeMsg)

		assert.Nil(t, err, "no error")
		assert.Equal(t, 1, len(kMessages))
		assert.Equal(t, "TestUser", string(kMessages[0].Value))
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
		changeMsg.Payload.Op = "u"

		kMessages, err := filterChange(changeMsg)

		assert.Nil(t, err, "no error")
		assert.Nil(t, kMessages, "nil output")
	})
}
