package definitions

import (
	"context"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/itchyny/gojq"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/parser/evaluation"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/pkg/utils"
	"github.com/blackstork-io/fabric/plugin"
)

type ParsedPlugin struct {
	PluginName string
	BlockName  string
	Meta       *MetaBlock
	Config     evaluation.Configuration
	Invocation evaluation.Invocation
}

func (pe *ParsedPlugin) GetBlockInvocation() *evaluation.BlockInvocation {
	res, ok := pe.Invocation.(*evaluation.BlockInvocation)
	if !ok {
		panic("This Plugin does not store a BlockInvocation!")
	}
	return res
}

type (
	ParsedContent ParsedPlugin
	ParsedData    ParsedPlugin
)

func (c *ParsedContent) Name() string {
	return c.PluginName
}

// Render implements Renderable.
func (c *ParsedContent) Render(ctx context.Context, caller evaluation.ContentCaller, dataCtx evaluation.DataContext, result *evaluation.Result, contentID uint32) (diags diagnostics.Diag) {
	if c.Meta != nil {
		dataCtx.Set(BlockKindContent, plugin.ConvMapData{
			BlockKindMeta: c.Meta.AsJQData(),
		})
	} else {
		dataCtx.Delete(BlockKindContent)
	}
	diags.Extend(c.EvalQuery(&dataCtx))
	// TODO: #28 #29
	if diags.HasErrors() {
		return
	}

	resultStr, diag := caller.CallContent(ctx, c.PluginName, c.Config, c.Invocation, dataCtx.AsJQData().(plugin.MapData))
	if diags.Extend(diag) || resultStr == nil {
		// XXX: What to do if we have errors while executing content blocks?
		// just skipping the value for now...
		return
	}
	if resultStr.Location == nil {
		resultStr.Location = &plugin.Location{
			Index: contentID,
		}
	}
	err := result.Add(resultStr.Content, resultStr.Location)
	if err != nil {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Failed to add content to the result",
			Detail:   err.Error(),
		})
	}
	return
}

func (c *ParsedContent) EvalQuery(dataCtx *evaluation.DataContext) (diags diagnostics.Diag) {
	body := c.Invocation.GetBody()
	attr, found := body.Attributes["query"]
	if !found {
		return
	}
	val, newBody, dgs := hcldec.PartialDecode(body, &hcldec.ObjectSpec{
		"query": &hcldec.AttrSpec{
			Name:     "query",
			Type:     cty.String,
			Required: true,
		},
	}, nil)
	c.Invocation.SetBody(utils.ToHclsyntaxBody(newBody))
	if diags.ExtendHcl(dgs) {
		return
	}
	query := val.GetAttr("query").AsString()

	dataCtx.Set("query", plugin.StringData(query))
	queryResult, err := runQuery(query, dataCtx)
	if err != nil {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Failed to execute the query",
			Detail:   err.Error(),
			Subject:  &attr.SrcRange,
		})
		return
	}
	dataCtx.Set("query_result", queryResult)
	return
}

func runQuery(query string, dataCtx *evaluation.DataContext) (result plugin.Data, err error) {
	jqQuery, err := gojq.Parse(query)
	if err != nil {
		err = fmt.Errorf("failed to parse the query: %w", err)
		return
	}

	code, err := gojq.Compile(jqQuery)
	if err != nil {
		err = fmt.Errorf("failed to compile the query: %w", err)
		return
	}
	res, hasResult := code.Run(dataCtx.Any()).Next()
	if !hasResult {
		return
	}
	result, err = plugin.ParseDataAny(res)
	if err != nil {
		err = fmt.Errorf("incorrect query result type: %w", err)
	}
	return
}

type Renderable interface {
	Name() string
	Render(ctx context.Context, caller evaluation.ContentCaller, dataCtx evaluation.DataContext, result *evaluation.Result, contentID uint32) diagnostics.Diag
}
