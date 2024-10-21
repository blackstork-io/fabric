package eval

import (
	"context"
	"log/slog"
	"maps"
	"slices"

	"github.com/hashicorp/hcl/v2"

	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/plugindata"
)

type Document struct {
	Source        *definitions.Document
	Meta          *definitions.MetaBlock
	Vars          *definitions.ParsedVars
	RequiredVars  []string
	DataBlocks    []*PluginDataAction
	ContentBlocks []*Content
	PublishBlocks []*PluginPublishAction
}

func (doc *Document) FetchData(ctx context.Context) (plugindata.Data, diagnostics.Diag) {
	logger := *slog.Default()
	logger.DebugContext(ctx, "Fetching data for the document template")
	result := make(plugindata.Map)
	diags := diagnostics.Diag{}
	for _, block := range doc.DataBlocks {
		var dsMap plugindata.Map
		found, ok := result[block.PluginName]
		if ok {
			dsMap = found.(plugindata.Map)
		} else {
			dsMap = make(plugindata.Map)
			result[block.PluginName] = dsMap
		}
		if _, found := dsMap[block.BlockName]; found {
			diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagWarning,
				Summary:  "Data conflict",
				Detail:   "Result of this block overwrites results from the previous invocation.",
				Subject:  &block.SrcRange,
			})
		}
		dsMap[block.BlockName], diags = block.FetchData(ctx)
		if diags.HasErrors() {
			return nil, diags
		}
	}
	return result, diags
}

func filterChildrenByTags(children []*Content, requiredTags []string) []*Content {
	return slices.DeleteFunc(children, func(child *Content) bool {
		switch {
		case child.Plugin != nil:
			return !child.Plugin.Meta.MatchesTags(requiredTags)
		case child.Section != nil:
			if child.Section.meta.MatchesTags(requiredTags) {
				return false
			}
			child.Section.children = filterChildrenByTags(child.Section.children, requiredTags)
			return len(child.Section.children) == 0
		}
		return false
	})
}

func (doc *Document) RenderContent(ctx context.Context, docDataCtx plugindata.Map, requiredTags []string) (*plugin.ContentSection, plugindata.Data, diagnostics.Diag) {
	logger := *slog.Default()
	logger.DebugContext(ctx, "Fetching data for the document template")
	data, diags := doc.FetchData(ctx)
	if diags.HasErrors() {
		return nil, nil, diags
	}
	docData := plugindata.Map{}
	if doc.Meta != nil {
		docData[definitions.BlockKindMeta] = doc.Meta.AsPluginData()
	}
	// static portion of the data context for this document
	// will never change, all changes are made to the clone of this map
	docDataCtx[definitions.BlockKindData] = data
	docDataCtx[definitions.BlockKindDocument] = docData

	diag := ApplyVars(ctx, doc.Vars, docDataCtx)

	if diags.Extend(diag) {
		return nil, nil, diags
	}

	// verify required vars
	if len(doc.RequiredVars) > 0 {
		diag = verifyRequiredVars(docDataCtx, doc.RequiredVars, doc.Source.Block)
		if diags.Extend(diag) {
			return nil, nil, diags
		}
	}

	// evaluate/expand dynamic blocks
	children, diag := UnwrapDynamicContent(ctx, doc.ContentBlocks, docDataCtx)
	if diags.Extend(diag) {
		return nil, nil, diags
	}
	// filter out content blocks that do not match tags
	if !doc.Meta.MatchesTags(requiredTags) {
		children = filterChildrenByTags(children, requiredTags)
	}

	result := plugin.NewSection(0)
	// create a position map for content blocks
	posMap := make(map[int]uint32, len(children))
	for i := range children {
		empty := new(plugin.ContentEmpty)
		result.Add(empty, nil)
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
	// execute content blocks based on the invocation order
	for _, idx := range invokeList {
		// clone the data context for each content block
		dataCtx := maps.Clone(docDataCtx)
		// set the current content to the data context
		dataCtx[definitions.BlockKindDocument].(plugindata.Map)[definitions.BlockKindContent] = result.AsData()
		// TODO: if section, set section

		// execute the content block
		diag := children[idx].RenderContent(ctx, dataCtx, result, result, posMap[idx])
		if diags.Extend(diag) {
			return nil, nil, diags
		}
	}
	// compact the content tree to remove empty content nodes
	result.Compact()
	return result, docDataCtx, diags
}

func (doc *Document) Publish(ctx context.Context, content plugin.Content, data plugindata.Data, documentName string) diagnostics.Diag {
	logger := *slog.Default()
	logger.DebugContext(ctx, "Fetching data for the document template")
	docData := plugindata.Map{
		definitions.BlockKindContent: content.AsData(),
	}
	if doc.Meta != nil {
		docData[definitions.BlockKindMeta] = doc.Meta.AsPluginData()
	}
	dataCtx := plugindata.Map{
		definitions.BlockKindData:     data,
		definitions.BlockKindDocument: docData,
	}
	var diags diagnostics.Diag
	for _, block := range doc.PublishBlocks {
		diag := block.Publish(ctx, dataCtx, documentName)
		if diag != nil {
			diags.Extend(diag)
		}
	}
	return diags
}

func LoadDocument(ctx context.Context, plugins Plugins, node *definitions.ParsedDocument) (_ *Document, diags diagnostics.Diag) {
	block := Document{
		Source:       node.Source,
		Meta:         node.Meta,
		Vars:         node.Vars,
		RequiredVars: node.RequiredVars,
	}
	dataNames := make(map[[2]string]struct{})
	for _, child := range node.Data {
		decoded, diag := LoadDataAction(ctx, plugins, child)
		if diags.Extend(diag) {
			return nil, diags
		}
		key := [2]string{decoded.PluginName, decoded.BlockName}
		if _, found := dataNames[key]; found {
			diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagWarning,
				Summary:  "Data conflict",
				Detail:   "Data block with the same name already exists.",
				Subject:  &decoded.SrcRange,
			})
		}
		dataNames[key] = struct{}{}
		block.DataBlocks = append(block.DataBlocks, decoded)
	}
	for _, child := range node.Content {
		decoded, diag := LoadContent(ctx, plugins, child)
		if diags.Extend(diag) {
			return nil, diags
		}
		block.ContentBlocks = append(block.ContentBlocks, decoded)
	}
	for _, child := range node.Publish {
		decoded, diag := LoadPluginPublishAction(ctx, plugins, child)
		if diags.Extend(diag) {
			return nil, diags
		}
		block.PublishBlocks = append(block.PublishBlocks, decoded)
	}
	return &block, diags
}
