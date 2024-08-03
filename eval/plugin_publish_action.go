package eval

import (
	"context"
	"fmt"
	"slices"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/pkg/utils"
	"github.com/blackstork-io/fabric/plugin"
)

type PluginPublishAction struct {
	*PluginAction
	Publisher *plugin.Publisher
	Format    plugin.OutputFormat
}

func (block *PluginPublishAction) Publish(ctx context.Context, dataCtx plugin.MapData, documentName string) diagnostics.Diag {
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
	var cfg cty.Value
	if p.Config != nil && !p.Config.IsEmpty() {
		cfg, diags = node.Config.ParseConfig(ctx, p.Config)
		if diags.HasErrors() {
			return nil, diags
		}
	} else if (p.Config == nil || p.Config.IsEmpty()) && node.Config.Exists() {
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
	body := node.Invocation.GetBody()
	var format plugin.OutputFormat
	if attr, found := body.Attributes["format"]; found {
		value, newBody, stdDiag := hcldec.PartialDecode(body, &hcldec.ObjectSpec{
			"format": &hcldec.AttrSpec{
				Name:     "format",
				Type:     cty.String,
				Required: true,
			},
		}, nil)
		if diags.Extend(stdDiag) {
			return
		}
		node.Invocation.SetBody(utils.ToHclsyntaxBody(newBody))
		formatStr := value.GetAttr("format").AsString()
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
		if !slices.Contains(p.AllowedFormats, format) {
			diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid format",
				Detail:   fmt.Sprintf("'%s' is not allowed for this publisher", format.String()),
				Subject:  &attr.SrcRange,
			})
			return
		}
	}
	var args cty.Value
	args, diag := node.Invocation.ParseInvocation(ctx, p.Args)
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
