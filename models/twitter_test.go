package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTwitterUserList(t *testing.T) {
	list := &TwitterUserList{}

	slice1 := []*TwitterUser{
		&TwitterUser{
			Username: "muesli",
		},
		&TwitterUser{
			Username: "smag",
		},
	}

	slice2 := []*TwitterUser{
		&TwitterUser{
			Username: "muesli",
		},
		&TwitterUser{
			Username: "smag",
		},
		&TwitterUser{
			Username: "franz",
		},
	}

	t.Run("create new TwitterUserList slice from sub-slices", func(t *testing.T) {
		list = NewTwitterUserList(slice1, slice2)
		assert.Equal(t, len(*list), 5, "should contain content of both slices with duplicates")
	})

	t.Run("remove all duplicated TwitterUser elements by username", func(t *testing.T) {
		list.RemoveDuplicates()
		assert.Equal(t, len(*list), 3, "should only include unique elements")
	})
}
