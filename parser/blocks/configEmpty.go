package blocks

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"

	"github.com/blackstork-io/fabric/parser/evaluation"
)

func NewEmptyConfig(MissingItemRange hcl.Range) *ConfigEmpty {
	return &ConfigEmpty{
		body: hclsyntax.Body{
			SrcRange: MissingItemRange,
			EndRange: hcl.Range{
				Filename: MissingItemRange.Filename,
				Start:    MissingItemRange.End,
				End:      MissingItemRange.End,
			},
		},
	}
}

// Empty config, storing the range of the original block
type ConfigEmpty struct {
	body hclsyntax.Body
}

// GetBody implements evaluation.Configuration.
func (c *ConfigEmpty) GetBody() hcl.Body {
	return &c.body
}

// Exists implements evaluation.Configuration.
func (c *ConfigEmpty) Exists() bool {
	return false
}

// Range implements Configuration.
func (c *ConfigEmpty) Range() hcl.Range {
	return c.body.SrcRange
}

var _ evaluation.Configuration = (*ConfigEmpty)(nil)
