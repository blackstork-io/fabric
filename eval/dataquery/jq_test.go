package dataquery

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/pkg/encapsulator"
)

func TestJqIsQuery(t *testing.T) {
	assert.True(t, encapsulator.Compatible(JqQueryType, definitions.QueryType))
}
