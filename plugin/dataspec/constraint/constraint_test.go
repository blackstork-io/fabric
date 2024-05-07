package constraint

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// lock constraints to their expected values
// changing these would be backwards incompatible (for protobuf)
func TestConstraints(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(1, Required)
	assert.Equal(2, NonNull)
	assert.Equal(4, NonEmpty)
	assert.Equal(8, TrimSpace)
	assert.Equal(16, Integer)
}
