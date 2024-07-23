package builtin

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/dataspec/constraint"
)

func makeTextContentProvider() *plugin.ContentProvider {
	return &plugin.ContentProvider{
		ContentFunc: genTextContent,
		Args: &dataspec.RootSpec{
			Attrs: []*dataspec.AttrSpec{
				{
					Name:        "value",
					Type:        cty.String,
					Constraints: constraint.RequiredNonNull,
					ExampleVal:  cty.StringVal("Hello world!"),
					Doc:         `A string to render. Can use go template syntax.`,
				},
			},
		},
		Doc: `Renders text`,
	}
}

func genTextContent(ctx context.Context, params *plugin.ProvideContentParams) (*plugin.ContentResult, diagnostics.Diag) {
	value := params.Args.GetAttrVal("value")
	if value.IsNull() {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse arguments",
			Detail:   "value is required",
		}}
	}

	text, err := genTextContentText(value.AsString(), params.DataContext)
	if err != nil {
		return nil, diagnostics.Diag{{
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
