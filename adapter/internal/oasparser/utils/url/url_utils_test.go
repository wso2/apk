package urlutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetURLType(t *testing.T) {
	type testItem struct {
		tls     bool
		urlType string
		message string
	}

	tests := []testItem{
		{
			tls:     true,
			urlType: "https",
			message: "Type should be https when tls is enabled",
		},
		{
			tls:     false,
			urlType: "http",
			message: "Type should be http when tls is disabled",
		},
	}

	for _, test := range tests {
		outputURLType := GetURLType(test.tls)
		assert.Equal(t, test.urlType, outputURLType, test.message)
	}
}
