package definitions

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/parser/evaluation"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
)

// Attribute referencing a configuration block (`config = path.to.config`).
type ConfigPtr struct {
	Cfg *Config
	Ptr *hcl.Attribute
}

// Parse implements ConfigurationObject.
func (c *ConfigPtr) Parse(spec hcldec.Spec) (val cty.Value, diags diagnostics.Diag) {
	return c.Cfg.Parse(spec)
}

// Range implements ConfigurationObject.
func (c *ConfigPtr) Range() hcl.Range {
	// Use the location of "config = *traversal*" for error reporting, not original config's Range
	return c.Ptr.Range
}

var _ evaluation.Configuration = (*ConfigPtr)(nil)
