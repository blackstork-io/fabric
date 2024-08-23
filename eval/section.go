package eval

import (
	"context"
	"fmt"
	"maps"
	"slices"

	"github.com/hashicorp/hcl/v2"

	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/plugindata"
)

type Section struct {
	meta         *definitions.MetaBlock
	children     []*Content
	vars         *definitions.ParsedVars
	source       *definitions.Section
	requiredVars []string
}

func (block *Section) RenderContent(ctx context.Context, dataCtx plugindata.Map, doc, parent *plugin.ContentSection, contentID uint32) (_ *plugin.ContentResult, diags diagnostics.Diag) {
	sectionData := plugindata.Map{}
	if block.meta != nil {
		sectionData[definitions.BlockKindMeta] = block.meta.AsPluginData()
	}
	section := new(plugin.ContentSection)
	if parent != nil {
		err := parent.Add(section, &plugin.Location{
			Index: contentID,
		})
		if err != nil {
			return nil, diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Failed to place content",
				Detail:   fmt.Sprintf("Failed to place content: %s", err),
			}}
		}
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
	dataCtx[definitions.BlockKindSection] = sectionData

	diag := ApplyVars(ctx, block.vars, dataCtx)
	if diags.Extend(diag) {
		return nil, diags
	}

	// verify required vars
	if len(block.requiredVars) > 0 {
		diag := verifyRequiredVars(dataCtx, block.requiredVars, block.source.Block)
		if diags.Extend(diag) {
			return nil, diags
		}
	}

	// execute content blocks based on the invocation order
	for _, idx := range invokeList {
		// update the session data (is propagated to dataCtx, maps are by-ref structures)
		sectionData[definitions.BlockKindContent] = section.AsData()

		// execute the content block
		_, diag := block.children[idx].RenderContent(ctx, maps.Clone(dataCtx), doc, section, posMap[idx])
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

func LoadSection(ctx context.Context, providers ContentProviders, node *definitions.ParsedSection) (_ *Section, diag diagnostics.Diag) {
	var diags diagnostics.Diag
	block := &Section{
		meta:         node.Meta,
		vars:         node.Vars,
		source:       node.Source,
		requiredVars: node.RequiredVars,
	}

	if node.Title != nil {
		title, diag := LoadContent(ctx, providers, node.Title)
		if diags.Extend(diag) {
			return nil, diags
		}
		block.children = append(block.children, title)

	}
	for _, child := range node.Content {
		decoded, diag := LoadContent(ctx, providers, child)
		if diags.Extend(diag) {
			return nil, diags
		}
		block.children = append(block.children, decoded)
	}
	return block, diags
}
