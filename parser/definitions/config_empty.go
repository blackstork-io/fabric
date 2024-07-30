package definitions

import (
	"context"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"

	"github.com/blackstork-io/fabric/parser/evaluation"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin/dataspec"
)

// Empty config, storing the range of the original block
type ConfigEmpty struct {
	Plugin *Plugin
}

// Exists implements evaluation.Configuration.
func (c *ConfigEmpty) Exists() bool {
	return false
}

// ParseConfig implements Configuration.
func (c *ConfigEmpty) ParseConfig(ctx context.Context, spec *dataspec.RootSpec) (val *dataspec.Block, diags diagnostics.Diag) {
	labels := make([]string, 1, len(c.Plugin.Block.Labels)+1)
	labels[0] = c.Plugin.Block.Type
	labels = append(labels, c.Plugin.Block.Labels...)
	labelRanges := make([]hcl.Range, 1, len(c.Plugin.Block.Labels)+1)
	labelRanges[0] = c.Plugin.Block.TypeRange
	labelRanges = append(labelRanges, c.Plugin.Block.LabelRanges...)
	if len(labels) >= 2 {
		// use the resolved name if it exists
		labels[1] = c.Plugin.Name()
	}

	emptyBody := hclsyntax.Block{
		Type:        "config",
		TypeRange:   c.Plugin.Block.TypeRange,
		Labels:      labels,
		LabelRanges: labelRanges,
		Body: &hclsyntax.Body{
			SrcRange: c.Plugin.Block.Body.MissingItemRange(),
			EndRange: c.Plugin.Block.Body.MissingItemRange(),
		},
		OpenBraceRange:  c.Plugin.Block.Body.MissingItemRange(),
		CloseBraceRange: c.Plugin.Block.Body.MissingItemRange(),
	}

	var diag diagnostics.Diag
	val, diag = dataspec.DecodeAndEvalBlock(ctx, &emptyBody, spec, nil)
	for _, d := range diag {
		d.Summary = fmt.Sprintf("Missing required configuration: %s", d.Summary)
	}
	return val, diagnostics.Diag(diag)
}

// Range implements Configuration.
func (c *ConfigEmpty) Range() hcl.Range {
	return c.Plugin.DefRange()
}

var _ evaluation.Configuration = (*ConfigEmpty)(nil)
