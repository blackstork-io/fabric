package builtin

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/pelletier/go-toml/v2"
	"github.com/zclconf/go-cty/cty"
	"gopkg.in/yaml.v3"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/dataspec/constraint"
	"github.com/blackstork-io/fabric/plugin/plugindata"
)

func makeFrontMatterContentProvider() *plugin.ContentProvider {
	return &plugin.ContentProvider{
		ContentFunc: genFrontMatterContent,
		Args: &dataspec.RootSpec{
			Attrs: []*dataspec.AttrSpec{
				{
					Name:       "format",
					Type:       cty.String,
					Doc:        `Format of the frontmatter.`,
					DefaultVal: cty.StringVal("yaml"),
					OneOf: []cty.Value{
						cty.StringVal("yaml"),
						cty.StringVal("toml"),
						cty.StringVal("json"),
					},
				},
				{
					Name:        "content",
					Type:        plugindata.EncapsulatedData.CtyType(),
					Doc:         `Arbitrary key-value map to be put in the frontmatter.`,
					Constraints: constraint.RequiredMeaningful,
					ExampleVal: cty.ObjectVal(map[string]cty.Value{
						"key": cty.StringVal("arbitrary value"),
						"key2": cty.MapVal(map[string]cty.Value{
							"can be nested": cty.NumberIntVal(42),
						}),
					}),
				},
			},
		},
		Doc: `Produces the frontmatter.`,
	}
}

func genFrontMatterContent(ctx context.Context, params *plugin.ProvideContentParams) (*plugin.ContentResult, diagnostics.Diag) {
	if err := validateFrontMatterContentTree(params.DataContext, params.ContentID); err != nil {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Error while validating frontmatter constraints",
			Detail:   err.Error(),
		}}
	}

	format := params.Args.GetAttrVal("format").AsString()

	data, err := plugindata.EncapsulatedData.FromCty(params.Args.GetAttrVal("content"))
	if err != nil {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse arguments",
			Detail:   err.Error(),
		}}
	}
	if data == nil || *data == nil {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse arguments",
			Detail:   "Content is nil",
		}}
	}
	m, ok := (*data).(plugindata.Map)
	if !ok {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse arguments",
			Detail:   fmt.Sprintf("Invalid frontmatter data type: %T. Map required.", data),
		}}
	}
	var result string
	switch format {
	case "yaml":
		result, err = renderYAMLFrontMatter(m)
	case "toml":
		result, err = renderTOMLFrontMatter(m)
	case "json":
		result, err = renderJSONFrontMatter(m)
	default:
		panic(fmt.Errorf("Unknown format type: %s", format))
	}
	if err != nil {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Failed to render frontmatter",
			Detail:   err.Error(),
		}}
	}
	return &plugin.ContentResult{
		Content: &plugin.ContentElement{
			Markdown: result,
		},
		Location: &plugin.Location{
			Index:  1,
			Effect: plugin.LocationEffectBefore,
		},
	}, nil
}

func validateFrontMatterContentTree(dataCtx plugindata.Map, contentID uint32) error {
	if dataCtx == nil {
		return fmt.Errorf("DataContext is empty")
	}
	document, _ := parseScope(dataCtx)
	if document == nil {
		return fmt.Errorf("frontmatter must be declared in the document")
	}
	if findDepth(document, contentID, 0) != 0 {
		return fmt.Errorf("frontmatter must be declared at the top-level of the document")
	}
	if countDeclarations(document, "frontmatter") > 0 {
		return fmt.Errorf("frontmatter already declared in the document")
	}
	return nil
}

func renderYAMLFrontMatter(m plugindata.Map) (string, error) {
	var buf strings.Builder
	buf.WriteString("---\n")
	err := yaml.NewEncoder(&buf).Encode(m)
	if err != nil {
		return "", err
	}
	buf.WriteString("---")
	return buf.String(), nil
}

func renderTOMLFrontMatter(m plugindata.Map) (string, error) {
	var buf strings.Builder
	buf.WriteString("+++\n")
	err := toml.NewEncoder(&buf).Encode(m)
	if err != nil {
		return "", err
	}
	buf.WriteString("+++")
	return buf.String(), nil
}

func renderJSONFrontMatter(m plugindata.Map) (string, error) {
	var buf strings.Builder
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	err := enc.Encode(m)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
