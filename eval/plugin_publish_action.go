package eval

import (
	"context"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/cmd/fabctx"
	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/pkg/utils"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/dataspec/constraint"
	"github.com/blackstork-io/fabric/plugin/plugindata"
)

type PluginPublishAction struct {
	*PluginAction
	Publisher *plugin.Publisher
	Format    plugin.OutputFormat
}

func (block *PluginPublishAction) Publish(ctx context.Context, dataCtx plugindata.Map, documentName string) diagnostics.Diag {
	return block.Publisher.Execute(ctx, &plugin.PublishParams{
		Config:       block.Config,
		Args:         block.Args,
		DataContext:  dataCtx,
		Format:       block.Format,
		DocumentName: documentName,
	})
}

func LoadPluginPublishAction(ctx context.Context, publishers Publishers, node *definitions.ParsedPlugin) (_ *PluginPublishAction, diags diagnostics.Diag) {
	p, ok := publishers.Publisher(node.PluginName)
	if !ok {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Missing publisher",
			Detail:   fmt.Sprintf("'%s' not found in any plugin", node.PluginName),
		}}
	}
	var cfg *dataspec.Block
	if p.Config != nil {
		cfg, diags = node.Config.ParseConfig(ctx, p.Config)
		if diags.HasErrors() {
			return nil, diags
		}
	} else if node.Config.Exists() {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagWarning,
			Summary:  "Publisher doesn't support configuration",
			Detail: fmt.Sprintf("Publisher '%s' does not support configuration, "+
				"but was provided with one. Remove it.", node.PluginName),
			Subject: node.Config.Range().Ptr(),
			Context: node.Invocation.Range().Ptr(),
		})
		return nil, diags
	}

	var format plugin.OutputFormat
	// XXX: So format is optional? Not including format in invocation doesn't validate it
	// anyway, this would change with the new AST
	if attr, found := utils.Pop(node.Invocation.Body.Attributes, "format"); found {
		val, diag := dataspec.DecodeAttr(&dataspec.AttrSpec{
			Name:        "format",
			Type:        cty.String,
			Constraints: constraint.RequiredMeaningful,
			OneOf: constraint.OneOf(utils.FnMap(p.AllowedFormats, func(f plugin.OutputFormat) cty.Value {
				return cty.StringVal(f.String())
			})),
		}, attr, fabctx.GetEvalContext(ctx))

		if diags.Extend(diag) {
			return
		}
		formatStr := val.AsString()
		switch formatStr {
		case plugin.OutputFormatMD.String():
			format = plugin.OutputFormatMD
		case plugin.OutputFormatHTML.String():
			format = plugin.OutputFormatHTML
		case plugin.OutputFormatPDF.String():
			format = plugin.OutputFormatPDF
		default:
			diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid format",
				Detail:   fmt.Sprintf("'%s' is not a valid format", formatStr),
				Subject:  &attr.SrcRange,
			})
			return
		}
	}

	args, diag := dataspec.DecodeAndEvalBlock(ctx, node.Invocation.Block, p.Args)
	if diags.Extend(diag) {
		return nil, diags
	}
	return &PluginPublishAction{
		PluginAction: &PluginAction{
			PluginName: node.PluginName,
			BlockName:  node.BlockName,
			Meta:       node.Meta,
			Config:     cfg,
			Args:       args,
		},
		Publisher: p,
		Format:    format,
	}, diags
}
