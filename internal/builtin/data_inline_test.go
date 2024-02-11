package builtin

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/plugin"
)

func Test_makeInlineDataSchema(t *testing.T) {
	schema := makeInlineDataSource()
	assert.Nil(t, schema.Config)
	assert.Nil(t, schema.Args)
	assert.NotNil(t, schema.DataFunc)
}

func Test_fetchInlineData(t *testing.T) {
	args := cty.ObjectVal(map[string]cty.Value{
		"foo": cty.StringVal("bar"),
		"baz": cty.NumberIntVal(1),
		"qux": cty.BoolVal(true),
		"quux": cty.ListVal([]cty.Value{
			cty.StringVal("corge"),
			cty.StringVal("grault"),
			cty.StringVal("garply"),
		}),
		"quuz": cty.ObjectVal(map[string]cty.Value{
			"garply": cty.StringVal("waldo"),
			"fred":   cty.NumberFloatVal(3.123),
			"plugh":  cty.BoolVal(false),
		}),
		"xyzzy": cty.NullVal(cty.String),
	})
	p := &plugin.Schema{
		DataSources: plugin.DataSources{
			"inline": makeInlineDataSource(),
		},
	}
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
