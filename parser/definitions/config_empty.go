package definitions

import (
	"context"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/cmd/fabctx"
	"github.com/blackstork-io/fabric/parser/evaluation"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin/dataspec"
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
func (c *ConfigEmpty) ParseConfig(ctx context.Context, spec dataspec.RootSpec) (val cty.Value, diags diagnostics.Diag) {
	emptyBody := &hclsyntax.Body{
		SrcRange: c.MissingItemRange,
		EndRange: hcl.Range{
			Filename: c.MissingItemRange.Filename,
			Start:    c.MissingItemRange.End,
			End:      c.MissingItemRange.End,
		},
	}

	var diag diagnostics.Diag
	val, diag = dataspec.Decode(emptyBody, spec, fabctx.GetEvalContext(ctx))
	for _, d := range diag {
		d.Summary = fmt.Sprintf("Missing required configuration: %s", d.Summary)
	}
	return val, diagnostics.Diag(diag)
}

// Range implements Configuration.
func (c *ConfigEmpty) Range() hcl.Range {
	return c.MissingItemRange
}

var _ evaluation.Configuration = (*ConfigEmpty)(nil)
