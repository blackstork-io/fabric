package blocks

import (
	"github.com/hashicorp/hcl/v2"
)

// Attribute referencing a configuration block (`config = path.to.config`).
type ConfigPtr struct {
	*Config
	Rng hcl.Range
}

// Exists implements evaluation.Configuration.
func (c *ConfigPtr) Exists() bool {
	return c != nil && c.Config.Exists()
}

// Range implements Configuration.
func (c *ConfigPtr) Range() hcl.Range {
	// Use the location of "config = *traversal*" for error reporting, not original config's Range
	return c.Rng
}
