package eval

import (
	"context"
	"maps"

	"github.com/blackstork-io/fabric/cmd/fabctx"
	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/dataspec/deferred"
	"github.com/blackstork-io/fabric/plugin/plugindata"
)

type Section struct {
	meta         *definitions.MetaBlock
	children     []*Content
	vars         *definitions.ParsedVars
	source       *definitions.Section
	requiredVars []string
	isIncluded   *dataspec.Attr
}

func (block *Section) RenderContent(ctx context.Context, dataCtx plugindata.Map) (res plugin.Content, diags diagnostics.Diag) {
	diag := ApplyVars(ctx, block.vars, dataCtx)
	if diags.Extend(diag) {
		return
	}

	isIncluded, diag := dataspec.EvalAttr(ctx, block.isIncluded, dataCtx)
	if diags.Extend(diag) {
		return
	}
	if isIncluded.IsNull() || !plugindata.IsTruthy(*plugindata.Encapsulated.MustFromCty(isIncluded)) {
		return
	}

	// verify required vars
	if len(block.requiredVars) > 0 {
		diag := verifyRequiredVars(dataCtx, block.requiredVars, block.source.Block)
		if diags.Extend(diag) {
			return
		}
	}

	children, diag := UnwrapDynamicContent(ctx, block.children, dataCtx)
	if diags.Extend(diag) {
		return
	}

	section := plugin.NewSection(block.meta, len(children))
	dataCtx[definitions.BlockKindSection] = section.AsPluginData()

	// execute content blocks
	for _, child := range children {
		// execute the content block
		n, diag := child.RenderContent(ctx, maps.Clone(dataCtx))
		if diags.Extend(diag) {
			continue
		}
		section.AppendChild(n)
	}
	res = section
	return
}

func LoadSection(ctx context.Context, providers ContentProviders, node *definitions.ParsedSection) (_ *Section, diags diagnostics.Diag) {
	block := &Section{
		meta:         node.Meta,
		vars:         node.Vars,
		source:       node.Source,
		requiredVars: node.RequiredVars,
	}
	var diag diagnostics.Diag
	isIncluded := node.IsIncluded
	if isIncluded == nil {
		isIncluded = defaultIsIncluded(node.Source.Block.DefRange())
	}

	block.isIncluded, diag = dataspec.DecodeAttr(
		fabctx.GetEvalContext(deferred.WithQueryFuncs(ctx)),
		isIncluded,
		isIncludedSpec,
	)
	if diags.Extend(diag) {
		return
	}

	if node.Title != nil {
		title, diag := LoadContent(ctx, providers, node.Title)
		if diags.Extend(diag) {
			return
		}
		block.children = append(block.children, title)

	}
	for _, child := range node.Content {
		decoded, diag := LoadContent(ctx, providers, child)
		if diags.Extend(diag) {
			return
		}
		block.children = append(block.children, decoded)
	}
	return block, diags
}
