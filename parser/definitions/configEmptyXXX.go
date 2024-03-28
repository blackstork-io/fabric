package definitions

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/parser/evaluation"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
)

// Empty config, storing the range of the original block
type ConfigEmpty struct {
	MissingItemRange hcl.Range
}

// Exists implements evaluation.Configuration.
func (c *ConfigEmpty) Exists() bool {
	return false
}

// ParseConfig implements Configuration.
func (c *ConfigEmpty) ParseConfig(spec hcldec.Spec) (val cty.Value, diags diagnostics.Diag) {
	emptyBody := &hclsyntax.Body{
		SrcRange: c.MissingItemRange,
		EndRange: hcl.Range{
			Filename: c.MissingItemRange.Filename,
			Start:    c.MissingItemRange.End,
			End:      c.MissingItemRange.End,
		},
	}

	var diag hcl.Diagnostics
	val, diag = hcldec.Decode(emptyBody, spec, nil)
	for _, d := range diag {
		d.Summary = "Missing required configuration: " + d.Summary
	}
	return val, diagnostics.Diag(diag)
}

// Range implements Configuration.
func (c *ConfigEmpty) Range() hcl.Range {
	return c.MissingItemRange
}

var _ evaluation.Configuration = (*ConfigEmpty)(nil)
