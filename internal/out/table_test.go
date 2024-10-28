package out

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_squashTo(t *testing.T) {
	tests := []struct {
		in       string
		max      int
		expected string
	}{
		{in: "test long string", max: 11, expected: "test ..ring"},
		{in: "test long string", max: 10, expected: "test..ring"},
	}

	for _, test := range tests {
		out := squashTo(test.in, test.max)
		assert.Equal(t, test.expected, out)
		assert.Equal(t, test.max, len(out))
	}
}
