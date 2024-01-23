package shadow_test

import (
	"testing"

	"github.com/infinimesh/infinimesh/pkg/shadow"
	pb "github.com/infinimesh/proto/shadow"
	"github.com/stretchr/testify/assert"
)

func TestKey(t *testing.T) {
	cases := []struct {
		device   string
		key      pb.StateKey
		expected string
	}{
		{"device", pb.StateKey_CONNECTION, "device:connection"},
		{"device", pb.StateKey_DESIRED, "device:desired"},
		{"device", pb.StateKey_REPORTED, "device:reported"},
		{"device", pb.StateKey(100), "device:garbage"},
	}

	for _, c := range cases {
		actual := shadow.Key(c.device, c.key)
		assert.Equal(t, c.expected, actual)
	}
}
