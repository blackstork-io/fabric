package parser_test

import (
	"testing"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/stretchr/testify/assert"

	parser_mocks "github.com/blackstork-io/fabric/mocks/parser"
	"github.com/blackstork-io/fabric/parser"
)

func TestAddIfMissing(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	m := map[string]*parser_mocks.FabricBlock{}

	m1 := parser_mocks.NewFabricBlock(t)
	m1.EXPECT().Block().Return(&hclsyntax.Block{})

	diag := parser.AddIfMissing(m, "key_1", m1)
	assert.Empty(diag)
	assert.Same(m1, m["key_1"])

	m2 := parser_mocks.NewFabricBlock(t)

	diag = parser.AddIfMissing(m, "key_2", m2)
	assert.Empty(diag)
	assert.Same(m1, m["key_1"])
	assert.Same(m2, m["key_2"])

	m3 := parser_mocks.NewFabricBlock(t)
	m3.EXPECT().Block().Return(&hclsyntax.Block{}).Once()

	diag = parser.AddIfMissing(m, "key_1", m3)
	assert.NotEmpty(diag)
	assert.Same(m1, m["key_1"])
	assert.Same(m2, m["key_2"])
}
