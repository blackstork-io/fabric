package definitions

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"

	"github.com/blackstork-io/fabric/pkg/utils"
)

func validateBlockName(block *hclsyntax.Block, idx int, required bool) *hcl.Diagnostic {
	if idx >= len(block.Labels) {
		if required {
			return &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Missing block name",
				Detail:   "Block name was not specified",
				Subject:  block.DefRange().Ptr(),
			}
		}
		return nil
	}

	if !hclsyntax.ValidIdentifier(block.Labels[idx]) {
		return &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid block name",
			Detail: fmt.Sprintf(
				"Block name '%s' is an invalid identifier",
				block.Labels[idx],
			),
			Subject: block.LabelRanges[idx].Ptr(),
			Context: block.DefRange().Ptr(),
		}
	}
	return nil
}

func validatePluginKind(block *hclsyntax.Block, kind string, kindRange hcl.Range) *hcl.Diagnostic {
	switch kind {
	case BlockKindContent, BlockKindData:
		return nil
	default:
		return &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid plugin kind",
			Detail: fmt.Sprintf(
				"Unknown plugin kind '%s', valid plugin kinds are '%s' and '%s'",
				kind, BlockKindContent, BlockKindData,
			),
			Subject: kindRange.Ptr(),
			Context: block.DefRange().Ptr(),
		}
	}
}

func validatePluginKindLabel(block *hclsyntax.Block, idx int) *hcl.Diagnostic {
	if idx >= len(block.Labels) {
		return &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Missing plugin kind",
			Detail:   "Plugin kind was not specified",
			Subject:  block.DefRange().Ptr(),
		}
	}

	return validatePluginKind(block, block.Labels[idx], block.LabelRanges[idx])
}

func validatePluginName(block *hclsyntax.Block, idx int) *hcl.Diagnostic {
	if idx >= len(block.Labels) {
		return &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Missing plugin name",
			Detail:   "Plugin name was not specified",
			Subject:  block.DefRange().Ptr(),
		}
	}
	return nil
}

// TODO: Is it a hard error or warning?
func validateLabelsLength(block *hclsyntax.Block, maxLabels int, labelUsage string) *hcl.Diagnostic {
	if len(block.Labels) > maxLabels {
		if labelUsage != "" {
			labelUsage = fmt.Sprintf("%s %s", block.Type, labelUsage)
		} else {
			labelUsage = block.Type
		}
		return &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("Invalid %s block", block.Type),
			Detail:   fmt.Sprintf("Too many lables, usage: '%s'", labelUsage),
			Subject:  hcl.RangeBetween(block.LabelRanges[maxLabels], block.LabelRanges[len(block.LabelRanges)-1]).Ptr(),
			Context:  block.DefRange().Ptr(),
		}
	}
	return nil
}

func NewNestingDiag(what string, block *hclsyntax.Block, body *hclsyntax.Body, validChildren []string) *hcl.Diagnostic {
	return &hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  "Invalid block type",
		Detail: fmt.Sprintf(
			"%s can't contain '%s' block, only %s",
			what,
			block.Type,
			utils.JoinSurround(", ", "'", validChildren...),
		),
		Subject: block.Range().Ptr(),
		Context: body.Range().Ptr(),
	}
}
