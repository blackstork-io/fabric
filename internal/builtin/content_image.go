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

	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/dataspec/constraint"
)

func makeImageContentProvider() *plugin.ContentProvider {
	return &plugin.ContentProvider{
		ContentFunc: genImageContent,
		Args: dataspec.ObjectSpec{
			&dataspec.AttrSpec{
				Name:        "src",
				Type:        cty.String,
				Constraints: constraint.RequiredMeaningfull,
				ExampleVal:  cty.StringVal("https://example.com/img.png"),
			},
			&dataspec.AttrSpec{
				Name:       "alt",
				Type:       cty.String,
				ExampleVal: cty.StringVal("Text description of the image"),
				// Not using empty string as DefaultVal here for semantical meaning
			},
		},
		Doc: "Returns an image tag",
	}
}

func genImageContent(ctx context.Context, params *plugin.ProvideContentParams) (*plugin.ContentResult, hcl.Diagnostics) {
	src := params.Args.GetAttr("src")
	if src.IsNull() || src.AsString() == "" {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse arguments",
			Detail:   "src is required",
		}}
	}
	alt := params.Args.GetAttr("alt")
	if alt.IsNull() {
		alt = cty.StringVal("")
	}

	srcStr, err := renderAsTemplate("src", src.AsString(), params.DataContext)
	if err != nil {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to render src value as a template",
			Detail:   err.Error(),
		}}
	}

	altStr, err := renderAsTemplate("alt", alt.AsString(), params.DataContext)
	if err != nil {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to render alt value as a template",
			Detail:   err.Error(),
		}}
	}

	// Make sure there are no line breaks in the values
	srcStr = strings.TrimSpace(strings.ReplaceAll(srcStr, "\n", ""))
	altStr = strings.TrimSpace(strings.ReplaceAll(altStr, "\n", ""))
	return &plugin.ContentResult{
		Content: &plugin.ContentElement{
			Markdown: fmt.Sprintf("![%s](%s)", altStr, srcStr),
		},
	}, nil
}

func renderAsTemplate(name string, value string, datactx plugin.MapData) (string, error) {

	if value == "" {
		return "", nil
	}

	tmpl, err := template.New(name).Funcs(sprig.FuncMap()).Parse(value)

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
