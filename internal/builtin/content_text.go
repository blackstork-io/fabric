package builtin

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/plugin"
)

func makeTextContentProvider() *plugin.ContentProvider {
	return &plugin.ContentProvider{
		ContentFunc: genTextContent,
		Args: hcldec.ObjectSpec{
			"value": &hcldec.AttrSpec{
				Name:     "value",
				Type:     cty.String,
				Required: true,
			},
		},
	}
}

func genTextContent(ctx context.Context, params *plugin.ProvideContentParams) (*plugin.ContentResult, hcl.Diagnostics) {
	value := params.Args.GetAttr("value")
	if value.IsNull() {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse arguments",
			Detail:   "value is required",
		}}
	}

	text, err := genTextContentText(value.AsString(), params.DataContext)
	if err != nil {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to render text",
			Detail:   err.Error(),
		}}
	}
	return &plugin.ContentResult{
		Content: &plugin.ContentElement{
			Markdown: text,
		},
	}, nil
}

func genTextContentText(text string, datactx plugin.MapData) (string, error) {
	tmpl, err := template.New("text").Funcs(sprig.FuncMap()).Parse(text)
	if err != nil {
		return "", fmt.Errorf("failed to parse text template: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, datactx.Any())
	if err != nil {
		return "", fmt.Errorf("failed to execute text template: %w", err)
	}
	return strings.TrimSpace(buf.String()), nil
}
