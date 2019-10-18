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

	// Create functionality
	list = NewTwitterUserList(slice1, slice2)
	assert.Equal(t, len(*list), 5, "should contain content of both slices with duplicates")

	// RemoveDuplicates functionality
	list.RemoveDuplicates()
	assert.Equal(t, len(*list), 3, "should only include unique elements")
}
