package parser

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
)

type DefinedBlocks struct {
	Config    map[Key]*Config
	Documents map[string]*Document
	Sections  map[string]*Section
	Plugins   map[Key]*Plugin
}

func NewDefinedBlocks() *DefinedBlocks {
	return &DefinedBlocks{
		Config:    map[Key]*Config{},
		Documents: map[string]*Document{},
		Sections:  map[string]*Section{},
		Plugins:   map[Key]*Plugin{},
	}
}

func (db *DefinedBlocks) Merge(other *DefinedBlocks) (diags diagnostics.Diag) {
	for k, v := range other.Config {
		diags.Append(AddIfMissing(db.Config, k, v))
	}
	for k, v := range other.Documents {
		diags.Append(AddIfMissing(db.Documents, k, v))
	}
	for k, v := range other.Sections {
		diags.Append(AddIfMissing(db.Sections, k, v))
	}
	for k, v := range other.Plugins {
		diags.Append(AddIfMissing(db.Plugins, k, v))
	}
	return
}

func AddIfMissing[M ~map[K]V, K comparable, V FabricBlock](m M, key K, newBlock V) *hcl.Diagnostic {
	if origBlock, found := m[key]; found {
		kind := origBlock.Block().Type
		origDefRange := origBlock.Block().DefRange()
		return &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("Duplicate '%s' declaration", kind),
			Detail:   fmt.Sprintf("'%s' with the same name originally defined at %s:%d", kind, origDefRange.Filename, origDefRange.Start.Line),
			Subject:  newBlock.Block().DefRange().Ptr(),
		}
	}
	m[key] = newBlock
	return nil
}
