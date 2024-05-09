package builtin

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/blackstork-io/fabric/pkg/diagnostics/diagtest"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/plugintest"
)

func Test_makeInlineDataSchema(t *testing.T) {
	schema := makeInlineDataSource()
	assert.Nil(t, schema.Config)
	assert.NotNil(t, schema.Args)
	assert.NotNil(t, schema.DataFunc)
}

func Test_fetchInlineData(t *testing.T) {
	p := &plugin.Schema{
		DataSources: plugin.DataSources{
			"inline": makeInlineDataSource(),
		},
	}
	args := plugintest.DecodeAndAssert(t, p.DataSources["inline"].Args, `
        foo  = "bar"
        baz  = 1
        qux  = true
        quux = ["corge", "grault", "garply"]
        quuz = {
                garply = "waldo"
                fred   = 3.123
                plugh  = false
        }
        xyzzy = null
    `, diagtest.Asserts{})
	data, diags := p.RetrieveData(context.Background(), "inline", &plugin.RetrieveDataParams{
		Args: args,
	})
	assert.Empty(t, diags)
	assert.Equal(t, plugin.MapData{
		"foo": plugin.StringData("bar"),
		"baz": plugin.NumberData(1),
		"qux": plugin.BoolData(true),
		"quux": plugin.ListData{
			plugin.StringData("corge"),
			plugin.StringData("grault"),
			plugin.StringData("garply"),
		},
		"quuz": plugin.MapData{
			"garply": plugin.StringData("waldo"),
			"fred":   plugin.NumberData(3.123),
			"plugh":  plugin.BoolData(false),
		},
		"xyzzy": nil,
	}, data)
}
