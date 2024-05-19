package eval

import (
	"context"
	"maps"
	"slices"

	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
)

type Section struct {
	meta     *definitions.MetaBlock
	children []*Content
}

func (block *Section) RenderContent(ctx context.Context, dataCtx plugin.MapData, doc, parent *plugin.ContentSection, contentID uint32) (_ *plugin.ContentResult, diags diagnostics.Diag) {
	sectionData := plugin.MapData{}
	if block.meta != nil {
		sectionData[definitions.BlockKindMeta] = block.meta.AsJQData()
	}
	section := new(plugin.ContentSection)
	if parent != nil {
		parent.Add(section, nil)
	}
	// create a position map for content blocks
	posMap := make(map[int]uint32)
	for i := range block.children {
		empty := new(plugin.ContentEmpty)
		section.Add(empty, nil)
		posMap[i] = empty.ID()
	}
	// sort content blocks by invocation order
	invokeList := make([]int, 0, len(block.children))
	for i := range block.children {
		invokeList = append(invokeList, i)
	}
	slices.SortStableFunc(invokeList, func(a, b int) int {
		ao := block.children[a].InvocationOrder()
		bo := block.children[b].InvocationOrder()
		return ao.Weight() - bo.Weight()
	})
	// execute content blocks based on the invocation order
	for _, idx := range invokeList {
		// clone the data context for each content block
		dataCtx = maps.Clone(dataCtx)
		// set the current content to the data context
		sectionData[definitions.BlockKindContent] = section.AsData()
		dataCtx[definitions.BlockKindSection] = sectionData
		// execute the content block
		_, diag := block.children[idx].RenderContent(ctx, dataCtx, doc, section, posMap[idx])
		if diags.Extend(diag) {
			return nil, diags
		}
	}
	// compact the content tree to remove empty content nodes
	section.Compact()
	return &plugin.ContentResult{
		Content: section,
		Location: &plugin.Location{
			Index: contentID,
		},
	}, diags
}

func LoadSection(providers ContentProviders, node *definitions.ParsedSection) (_ *Section, diag diagnostics.Diag) {
	var diags diagnostics.Diag
	block := &Section{
		meta: node.Meta,
	}

	if node.Title != nil {
		title, diag := LoadContent(providers, node.Title)
		if diags.Extend(diag) {
			return nil, diags
		}
		block.children = append(block.children, title)

	}
	for _, child := range node.Content {
		decoded, diag := LoadContent(providers, child)
		if diags.Extend(diag) {
			return nil, diags
		}
		block.children = append(block.children, decoded)
	}
	return block, diags
}
