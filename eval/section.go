package eval

import (
	"context"
	"fmt"
	"maps"
	"slices"

	"github.com/hashicorp/hcl/v2"

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

func (block *Section) PrepareData(ctx context.Context, dataCtx plugindata.Map, doc, parent *plugin.ContentSection) (diags diagnostics.Diag) {
	sectionData := plugindata.Map{}
	if block.meta != nil {
		sectionData[definitions.BlockKindMeta] = block.meta.AsPluginData()
	}
	dataCtx[definitions.BlockKindSection] = sectionData
	diag := ApplyVars(ctx, block.vars, dataCtx)
	if diags.Extend(diag) {
		return
	}

	// verify required vars
	if len(block.requiredVars) > 0 {
		diag := verifyRequiredVars(dataCtx, block.requiredVars, block.source.Block)
		if diags.Extend(diag) {
			return
		}
	}

	return diags
}

func (block *Section) Unwrap(ctx context.Context, dataCtx plugindata.Map) (include bool, children []*Content, diags diagnostics.Diag) {
	isIncluded, diag := dataspec.EvalAttr(ctx, block.isIncluded, dataCtx)
	if diags.Extend(diag) {
		return
	}
	if isIncluded.IsNull() || !plugindata.IsTruthy(*plugindata.Encapsulated.MustFromCty(isIncluded)) {
		return
	}

	children, diag = UnwrapDynamicContent(ctx, block.children, dataCtx)
	if diags.Extend(diag) {
		return false, nil, diags
	}
	return true, children, diags
}

func (block *Section) RenderContent(ctx context.Context, dataCtx plugindata.Map, doc, parent *plugin.ContentSection, contentID uint32) (diags diagnostics.Diag) {
	sectionData := plugindata.Map{}
	if block.meta != nil {
		sectionData[definitions.BlockKindMeta] = block.meta.AsPluginData()
	}
	dataCtx[definitions.BlockKindSection] = sectionData
	section := new(plugin.ContentSection)
	if parent != nil {
		err := parent.Add(section, &plugin.Location{
			Index: contentID,
		})
		if err != nil {
			return diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Failed to place content",
				Detail:   fmt.Sprintf("Failed to place content: %s", err),
			}}
		}
	}

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

	children, diag := UnwrapDynamicContent(ctx, block.children, dataCtx)
	if diags.Extend(diag) {
		return
	}

	// create a position map for content blocks
	posMap := make(map[int]uint32)
	for i := range children {
		empty := new(plugin.ContentEmpty)
		section.Add(empty, nil)
		posMap[i] = empty.ID()
	}
	// sort content blocks by invocation order
	invokeList := make([]int, 0, len(children))
	for i := range children {
		invokeList = append(invokeList, i)
	}
	slices.SortStableFunc(invokeList, func(a, b int) int {
		ao := children[a].InvocationOrder()
		bo := children[b].InvocationOrder()
		return ao.Weight() - bo.Weight()
	})

	// verify required vars
	if len(block.requiredVars) > 0 {
		diag := verifyRequiredVars(dataCtx, block.requiredVars, block.source.Block)
		if diags.Extend(diag) {
			return
		}
	}

	// execute content blocks based on the invocation order
	for _, idx := range invokeList {
		// update the session data (is propagated to dataCtx, maps are by-ref structures)
		sectionData[definitions.BlockKindContent] = section.AsData()

		// execute the content block
		diag := children[idx].RenderContent(ctx, maps.Clone(dataCtx), doc, section, posMap[idx])
		if diags.Extend(diag) {
			return
		}
	}
	// compact the content tree to remove empty content nodes
	section.Compact()
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
