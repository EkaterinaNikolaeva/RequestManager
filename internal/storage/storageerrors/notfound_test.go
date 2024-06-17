package storageerrors

import (
	"testing"

	"github.com/tj/assert"
)

func TestNewNotFoundError(t *testing.T) {
	tests := map[string]struct {
		expected string
	}{
		"correct": {
			expected: "the object was not found in the storage",
		},
	}
	for testName, tc := range tests {
		t.Run(testName, func(t *testing.T) {
			err := NewNotFoundError()
			assert.Equal(t, tc.expected, err.Error())
		})

	}
}
