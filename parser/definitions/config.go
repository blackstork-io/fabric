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

// Configuration block.
type Config struct {
	*hcl.Block
	blockRange hcl.Range
	once       sync.Once
	value      cty.Value
}

// Exists implements evaluation.Configuration.
func (c *Config) Exists() bool {
	return c != nil
}

var _ evaluation.Configuration = (*Config)(nil)

// ParseConfig implements Configuration.
func (c *Config) ParseConfig(spec hcldec.Spec) (val cty.Value, diags diagnostics.Diag) {
	c.once.Do(func() {
		var diag hcl.Diagnostics
		c.value, diag = hcldec.Decode(c.Body, spec, evaluation.NewEvalContext())
		if diags.ExtendHcl(diag) {
			// don't let partially-decoded values live
			c.value = cty.NilVal
		}
	})
	val = c.value
	if val.IsNull() && diags == nil {
		diags.Append(diagnostics.RepeatedError)
	}
	return
}

// Range implements Configuration.
func (c *Config) Range() hcl.Range {
	return c.blockRange
}

func (c *Config) GetKey() *Key {
	var name string
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

func (c *Config) ApplicableTo(plugin *Plugin) bool {
	switch len(c.Block.Labels) {
	case 0:
		// anonymous config block
		return true
	case 2, 3:
		// named config block
		return plugin.Kind() == c.Block.Labels[0] && plugin.Name() == c.Block.Labels[1]
	default:
		panic("Invalid parsed config block")
	}
}

var _ FabricBlock = (*Config)(nil)

func (c *Config) GetHCLBlock() *hcl.Block {
	return c.Block
}

var ctyConfigType = capsuleTypeFor[Config]()

func (c *Config) CtyType() cty.Type {
	return ctyConfigType
}

func DefineConfig(block *hclsyntax.Block) (config *Config, diags diagnostics.Diag) {
	diags.Append(validatePluginKindLabel(block, 0))
	diags.Append(validatePluginName(block, 1))
	diags.Append(validateBlockName(block, 2, false))
	diags.Append(validateLabelsLength(block, 3, "plugin_kind plugin_name <block_name>"))

	if diags.HasErrors() {
		return
	}
	config = &Config{
		Block:      block.AsHCLBlock(),
		blockRange: block.Range(),
	}
	return
}

type ConfigResolver func(pluginKind, pluginName string) (config *Config)
