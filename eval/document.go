package eval

import (
	"context"
	"fmt"
	"log/slog"
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
	FormatBlocks  []*PluginFormatAction
}

func (doc *Document) FetchData(ctx context.Context) (plugindata.Data, diagnostics.Diag) {
	evaluator := makeAsyncDataEvaluator(ctx, doc, slog.Default())
	return evaluator.Execute()
}

func (doc *Document) FetchDataWithPath(ctx context.Context, path []string) (plugindata.Data, diagnostics.Diag) {
	evaluator := makeAsyncDataEvaluatorWithPath(ctx, doc, path, slog.Default())
	return evaluator.Execute()
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

func (doc *Document) RenderContent(
	ctx context.Context,
	docDataCtx plugindata.Map,
	requiredTags []string,
) (*plugin.ContentSection, plugindata.Data, diagnostics.Diag) {
	log := slog.Default()
	log.InfoContext(ctx, "Rendering document content", "document", doc.Source.Name)
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

	evaluator, diag := makeAsyncContentEvaluator(ctx, children, docDataCtx)
	if diags.Extend(diag) {
		return nil, nil, diags
	}

	result, diag := evaluator.Execute(docDataCtx)
	if diags.Extend(diag) {
		return nil, nil, diags
	}

	return result, docDataCtx, diags
}

func (doc *Document) Publish(
	ctx context.Context,
	content plugin.Content,
	data plugindata.Data,
	documentName string,
) diagnostics.Diag {
	log := *slog.Default()
	log.DebugContext(ctx, "Publishing a document")

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

	var formatDiags diagnostics.Diag

	requiredFormats := map[string]*PluginFormatAction{}
	for _, block := range doc.PublishBlocks {

		if block.Format == nil {
			continue
		}

		formatToPublish := *block.Format
		if _, ok := requiredFormats[formatToPublish]; !ok {

			var formatBlock *PluginFormatAction
			for _, fb := range doc.FormatBlocks {
				if fb.PluginName == formatToPublish {
					formatBlock = fb
					break
				}
			}

			if formatBlock == nil {
				formatDiags.Extend(diagnostics.Diag{{
					Severity: hcl.DiagError,
					Summary:  "Missing formatter requested by publisher",
					Detail: fmt.Sprintf(
						"'%s' formatter required by '%s' publisher not found",
						block.PluginName,
						formatToPublish,
					),
				}})
				continue
			}

			requiredFormats[formatToPublish] = formatBlock
		}
	}
	if formatDiags.HasErrors() {
		return formatDiags
	}

	var diags diagnostics.Diag
	for _, block := range doc.PublishBlocks {

		log.DebugContext(
			ctx, "Executing a publish block",
			"publisher", block.PluginName,
			"block_name", block.BlockName,
		)

		var formatter *PluginFormatAction
		if block.Format != nil {
			formatter = requiredFormats[*block.Format]
		}

		diag := block.Publish(ctx, dataCtx, documentName, formatter)
		if diag != nil {
			diags.Extend(diag)
		}
	}
	return diags
}

func LoadDocument(
	ctx context.Context,
	plugins Plugins,
	node *definitions.ParsedDocument,
) (_ *Document, diags diagnostics.Diag) {
	doc := Document{
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
		doc.DataBlocks = append(doc.DataBlocks, decoded)
	}
	for _, child := range node.Content {
		decoded, diag := LoadContent(ctx, plugins, child)
		if diags.Extend(diag) {
			return nil, diags
		}
		doc.ContentBlocks = append(doc.ContentBlocks, decoded)
	}

	for _, child := range node.Format {
		decoded, diag := LoadPluginFormatAction(ctx, plugins, child)
		if diags.Extend(diag) {
			return nil, diags
		}
		doc.FormatBlocks = append(doc.FormatBlocks, decoded)
	}

	// register MD format block if it's not defined
	mdIdx := slices.IndexFunc(doc.FormatBlocks, func(b *PluginFormatAction) bool {
		return b.Formatter.Format == "md"
	})
	if mdIdx == -1 {
		mdFormatAction, diag := CreatePluginFormatActionMD(ctx, plugins)
		doc.FormatBlocks = append(doc.FormatBlocks, mdFormatAction)
		if diags.Extend(diag) {
			return nil, diags
		}
	}

	for _, child := range node.Publish {
		decoded, diag := LoadPluginPublishAction(ctx, plugins, child)
		if diags.Extend(diag) {
			return nil, diags
		}
		doc.PublishBlocks = append(doc.PublishBlocks, decoded)
	}
	return &doc, diags
}
