package parsetree

import (
	"context"
	"fmt"
	"maps"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/itchyny/gojq"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/parser/plugincaller"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/pkg/utils"
	"github.com/blackstork-io/fabric/plugin"
)

type DocumentNode struct {
	meta         *definitions.MetaBlock
	contentNodes []Renderable
	dataNodes    []*definitions.ParsedPlugin
}

func (dn *DocumentNode) AddContent(content *definitions.ParsedPlugin) {
	dn.contentNodes = append(dn.contentNodes, ContentNode{plugin: content})
}

func (dn *DocumentNode) AddSection(section *definitions.ParsedSection) {
	dn.contentNodes = append(dn.contentNodes, SectionNode{section: section})
}

func (dn *DocumentNode) AddData(data *definitions.ParsedPlugin) {
	dn.dataNodes = append(dn.dataNodes, data)
}

func (dn *DocumentNode) AddMeta(meta *definitions.MetaBlock) {
	dn.meta = meta
}

// result has a shape map[plugin_name]map[block_name]plugin_result.
func (dn *DocumentNode) EvalData(ctx context.Context, caller plugincaller.DataCaller) (result plugin.MapData, diags diagnostics.Diag) {
	// TODO: can be parallel:
	// TODO: once again: meta blocks are not used

	result = plugin.MapData{}
	for _, node := range dn.dataNodes {
		res, diag := caller.CallData(
			ctx,
			node.PluginName,
			node.Config,
			node.Invocation,
		)
		if diags.Extend(diag) {
			continue
		}

		var pluginNameRes plugin.MapData
		if m, found := result[node.PluginName]; found {
			pluginNameRes = m.(plugin.MapData)
		} else {
			pluginNameRes = plugin.MapData{}
			result[node.PluginName] = pluginNameRes
		}
		pluginNameRes[node.BlockName] = res
	}
	return
}

type resultsList struct {
	ptr *[]string
}

func (d resultsList) AsJQ() plugin.Data {
	lst := *d.ptr
	dst := make([]plugin.Data, len(lst))
	for i, v := range lst {
		dst[i] = plugin.StringData(v)
	}
	return plugin.ListData(dst)
}

func (dn *DocumentNode) EvalContent(ctx context.Context, caller plugincaller.PluginCaller) (result []string, diags diagnostics.Diag) {
	dataResult, diags := dn.EvalData(ctx, caller)
	if diags.HasErrors() {
		return
	}

	document := plugin.ConvMapData{
		definitions.BlockKindContent: resultsList{ptr: &result},
	}
	if dn.meta != nil {
		document[definitions.BlockKindMeta] = dn.meta.AsJQ()
	}

	globalCtx := plugin.ConvMapData{
		definitions.BlockKindData:     dataResult,
		definitions.BlockKindDocument: document,
	}

	for _, content := range dn.contentNodes {
		localCtx := maps.Clone(globalCtx)
		diags.Extend(
			content.Render(ctx, caller, localCtx, &result),
		)
	}
	return
}

type ContentNode struct {
	plugin *definitions.ParsedPlugin
}

// Render implements Renderable.
func (c ContentNode) Render(ctx context.Context, caller plugincaller.ContentCaller, localCtx plugin.ConvMapData, result *[]string) (diags diagnostics.Diag) {
	if c.plugin.Meta != nil {
		localCtx[definitions.BlockKindContent] = plugin.ConvMapData{
			definitions.BlockKindMeta: c.plugin.Meta.AsJQ(),
		}
	}

	query, found, rng, diag := c.GetQuery()
	if !diags.Extend(diag) && found {
		diags.Extend(ExecuteQuery(query, rng, localCtx))
	}
	// TODO: #28 #29
	if diags.HasErrors() {
		return
	}

	resultStr, diag := caller.CallContent(ctx, c.plugin.PluginName, c.plugin.Config, c.plugin.Invocation, localCtx.AsJQ().(plugin.MapData))
	if diags.Extend(diag) {
		// XXX: What to do if we have errors while executing content blocks?
		// just skipping the value for now...
		return
	}
	*result = append(*result, resultStr)
	return
}

// GetQuery implements Queriable.
func (c ContentNode) GetQuery() (query string, found bool, rng *hcl.Range, diags diagnostics.Diag) {
	body := c.plugin.Invocation.GetBody()
	attr, found := body.Attributes["query"]
	if !found {
		return
	}
	rng = &attr.SrcRange
	val, newBody, dgs := hcldec.PartialDecode(body, &hcldec.ObjectSpec{
		"query": &hcldec.AttrSpec{
			Name:     "query",
			Type:     cty.String,
			Required: true,
		},
	}, nil)
	c.plugin.Invocation.SetBody(utils.ToHclsyntaxBody(newBody))
	if diags.ExtendHcl(dgs) {
		return
	}
	query = val.GetAttr("query").AsString()
	return
}

var _ Queriable = ContentNode{}

type SectionNode struct {
	section *definitions.ParsedSection
}

// Render implements Renderable.
func (s SectionNode) Render(ctx context.Context, caller plugincaller.ContentCaller, globalCtx plugin.ConvMapData, result *[]string) (diags diagnostics.Diag) {
	localCtx := maps.Clone(globalCtx)
	if s.section.Meta != nil {
		localCtx[definitions.BlockKindSection] = plugin.ConvMapData{
			definitions.BlockKindMeta: s.section.Meta.AsJQ(),
		}
	}
	if title := s.section.Title; title != nil {
		diags.Extend(ContentNode{plugin: title}.Render(ctx, caller, localCtx, result))
	}

	for _, content := range s.section.Content {
		switch contentT := content.(type) {
		case *definitions.ParsedPlugin:
			localLocalCtx := maps.Clone(localCtx)
			content := ContentNode{plugin: contentT}
			diags.Extend(
				content.Render(ctx, caller, localLocalCtx, result),
			)
		case *definitions.ParsedSection:
			diags.Extend(SectionNode{section: contentT}.Render(ctx, caller, localCtx, result))
		default:
			panic("must be exhaustive")
		}
	}
	return
}

func ExecuteQuery(query string, rng *hcl.Range, localCtx plugin.ConvMapData) (diags diagnostics.Diag) {
	localCtx["query"] = plugin.StringData(query)
	queryResult, err := runQuery(query, localCtx)
	if err != nil {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Failed to execute the query",
			Detail:   err.Error(),
			Subject:  rng,
		})
		return
	}
	localCtx["query_result"] = queryResult
	return
}

func runQuery(query string, dataCtx plugin.ConvMapData) (queryResult plugin.Data, err error) {
	q, err := gojq.Parse(query)
	if err != nil {
		err = fmt.Errorf("failed to parse the query: %w", err)
		return
	}

	code, err := gojq.Compile(q)
	if err != nil {
		err = fmt.Errorf("failed to compile the query: %w", err)
		return
	}
	res, hasResult := code.Run(dataCtx.Any()).Next()
	if hasResult {
		queryResult, err = plugin.ParseDataAny(res)
		if err != nil {
			err = fmt.Errorf("incorrect query result type: %w", err)
		}
	}
	return
}

type Renderable interface {
	Render(ctx context.Context, caller plugincaller.ContentCaller, localCtx plugin.ConvMapData, result *[]string) diagnostics.Diag
}

type Queriable interface {
	GetQuery() (query string, found bool, rng *hcl.Range, diags diagnostics.Diag)
}
