package parser

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"

	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
)

// Moving around the DefinedBlocks (for refs and config attrs)

func traversalFromExpr(expr hcl.Expression) (path []string, diags diagnostics.Diag) {
	// ignore diags, just checking if the val is null
	val, _ := expr.Value(nil)
	if val.IsNull() {
		// empty ref
		return
	}
	traversal, diag := hcl.AbsTraversalForExpr(expr)
	if diags.ExtendHcl(diag) {
		return
	}
	path = make([]string, len(traversal))
	for i, trav := range traversal {
		switch traverser := trav.(type) {
		case hcl.TraverseRoot:
			path[i] = traverser.Name
		case hcl.TraverseAttr:
			path[i] = traverser.Name
		default:
			diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid path",
				Detail:   "The path in the attribute can not contain this operation",
				Subject:  traverser.SourceRange().Ptr(),
			})
		}
	}
	if diag.HasErrors() {
		path = nil
	}
	return
}

func (db *DefinedBlocks) Traverse(expr hcl.Expression) (res any, diags diagnostics.Diag) {
	var found bool
	path, diags := traversalFromExpr(expr)
	if diags.HasErrors() {
		return
	}
	if len(path) == 0 {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid path",
			Detail:   "The path is empty",
			Subject:  expr.Range().Ptr(),
		})
		return
	}
	switch path[0] {
	case definitions.BlockKindConfig:
		if len(path) != 4 {
			diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid path",
				Detail:   "The config path should have format config.<plugin_kind>.<plugin_name>.<config_name>",
				Subject:  expr.Range().Ptr(),
			})
			return
		}
		res, found = db.Config[definitions.Key{
			PluginKind: path[1],
			PluginName: path[2],
			BlockName:  path[3],
		}]
	case definitions.BlockKindContent, definitions.BlockKindData:
		if len(path) != 3 {
			diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid path",
				Detail: fmt.Sprintf(
					"The %s path should have format %s.<plugin_kind>.<plugin_name>",
					path[0], path[0],
				),
				Subject: expr.Range().Ptr(),
			})
			return
		}
		res, found = db.Plugins[definitions.Key{
			PluginKind: path[0],
			PluginName: path[1],
			BlockName:  path[2],
		}]
	case definitions.BlockKindSection:
		switch len(path) {
		case 1, 2:
		default:
			diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid path",
				Detail:   "The section path should have format section.<optionally_ref>.<section_name>",
				Subject:  expr.Range().Ptr(),
			})
			return
		}
		res, found = db.Sections[path[len(path)-1]]
	default:
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid path",
			Detail:   fmt.Sprintf("Unknown path root '%s'", path[0]),
			Subject:  expr.Range().Ptr(),
		})
		return
	}
	if !found {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid path",
			Detail:   "Referenced item not found",
			Subject:  expr.Range().Ptr(),
		})
	}
	return
}
