package blocks

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/parser/blocks/internal/tree"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
)

// Configuration block.
type Config struct {
	tree.NodeSigil
	Block *hclsyntax.Block
}

var _ FabricBlock = &Config{}

var ctyConfigType = capsuleTypeFor[Config]()

func (c *Config) CtyType() cty.Type {
	return ctyConfigType
}

func (c *Config) AsCtyValue() cty.Value {
	return cty.CapsuleVal(ctyConfigType, c)
}

func (c *Config) HCLBlock() *hclsyntax.Block {
	return c.Block
}

func (*Config) FriendlyName() string {
	return "config block"
}

// Exists implements evaluation.Configuration.
func (c *Config) Exists() bool {
	return c != nil
}

// ParseConfig implements Configuration.
func (c *Config) GetBody() hcl.Body {
	return c.Block.Body
}

// Range implements Configuration.
func (c *Config) Range() hcl.Range {
	return c.Block.Range()
}

func (c *Config) ApplicabilityTest(pluginKind, pluginName string, applicationLoc *hcl.Range) *hcl.Diagnostic {
	switch len(c.Block.Labels) {
	case 0:
		// anonymous config block is always applicable
		return nil
	case 2, 3:
	default:
		panic("Invalid parsed config block")
	}
	// named config block
	if pluginKind == c.Block.Labels[0] && pluginName == c.Block.Labels[1] {
		return nil
	}
	return &hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  "Config not applicable to the given block",
		Detail: fmt.Sprintf(
			"Config for '%s.%s' was applied to '%s.%s'",
			c.Block.Labels[0],
			c.Block.Labels[1],
			pluginKind,
			pluginName,
		),
		Subject: applicationLoc,
	}
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
		Block: block,
	}
	return
}
