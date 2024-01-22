package definitions

import (
	"sync"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/parser/evaluation"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
)

// Configuration block
type Config struct {
	*hcl.Block
	BlockRange hcl.Range
	once       sync.Once
	val        cty.Value
}

var _ evaluation.Configuration = (*Config)(nil)

// Parse implements Configuration.
func (c *Config) Parse(spec hcldec.Spec) (val cty.Value, diags diagnostics.Diag) {
	c.once.Do(func() {
		var diag hcl.Diagnostics
		c.val, diag = hcldec.Decode(c.Body, spec, nil)
		if diags.ExtendHcl(diag) {
			// don't let partially-decoded values live
			c.val = cty.NilVal
		}
	})
	val = c.val
	if val.IsNull() && diags == nil {
		diags.Append(diagnostics.RepeatedError)
	}
	return
}

// Range implements ConfigurationObject.
func (c *Config) Range() hcl.Range {
	return c.BlockRange
}

func (c *Config) GetKey() *Key {
	name := ""
	switch len(c.Block.Labels) {
	case 0:
		// anonymous config block
		return nil
	case 3:
		// named config block
		name = c.Block.Labels[2]
		fallthrough
	case 2:
		// default config block
		return &Key{
			PluginKind: c.Block.Labels[0],
			PluginName: c.Block.Labels[1],
			BlockName:  name,
		}
	default:
		panic("Invalid parsed config block")
	}
}

var _ FabricBlock = (*Config)(nil)

func (c *Config) GetHCLBlock() *hcl.Block {
	return c.Block
}

func DefineConfig(block *hclsyntax.Block) (config *Config, diags diagnostics.Diag) {
	diags.Append(validatePluginKindLabel(block, 0))
	diags.Append(validatePluginName(block, 1))
	diags.Append(validateBlockName(block, 2, true))
	diags.Append(validateLabelsLength(block, 3, "plugin_kind plugin_name block_name"))

	if diags.HasErrors() {
		return
	}
	config = &Config{
		Block:      block.AsHCLBlock(),
		BlockRange: block.Range(),
	}
	return
}
