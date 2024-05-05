package builtin

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/pelletier/go-toml/v2"
	"github.com/zclconf/go-cty/cty"
	"gopkg.in/yaml.v3"

	"github.com/blackstork-io/fabric/pkg/utils"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
)

const (
	frontMatterQueryResultKey = "query_result"
)

var frontMatterAllowedFormats = []string{"yaml", "toml", "json"}

func makeFrontMatterContentProvider() *plugin.ContentProvider {
	return &plugin.ContentProvider{
		ContentFunc: genFrontMatterContent,
		Args: dataspec.ObjectSpec{
			&dataspec.AttrSpec{
				Name:       "format",
				Type:       cty.String,
				Doc:        `Format of the frontmatter. Must be one of ` + utils.JoinSurround(", ", `"`, frontMatterAllowedFormats...),
				DefaultVal: cty.StringVal("yaml"),
			},
			&dataspec.AttrSpec{
				Name: "content",
				Type: cty.DynamicPseudoType,
				Doc: `
				Arbitrary key-value map to be put in the frontmatter.

				NOTE: Data from "query_result" replaces this value if present`,
				DefaultVal: cty.NullVal(cty.DynamicPseudoType),
				ExampleVal: cty.ObjectVal(map[string]cty.Value{
					"key": cty.StringVal("arbitrary value"),
					"key2": cty.MapVal(map[string]cty.Value{
						"can be nested": cty.NumberIntVal(42),
					}),
				}),
			},
		},
		Doc: `Produces the frontmatter.`,
	}
}

func genFrontMatterContent(ctx context.Context, params *plugin.ProvideContentParams) (*plugin.ContentResult, hcl.Diagnostics) {
	if err := validateFrontMatterContentTree(params.DataContext, params.ContentID); err != nil {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Error while validating frontmatter constraints",
			Detail:   err.Error(),
		}}
	}
	format, m, err := parseFrontMatterArgs(params.Args, params.DataContext)
	if err != nil {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse arguments",
			Detail:   err.Error(),
		}}
	}
	result, err := renderFrontMatterContent(format, m)
	if err != nil {
		return nil, hcl.Diagnostics{{
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

func validateFrontMatterContentTree(datactx plugin.MapData, contentID uint32) error {
	if datactx == nil {
		return fmt.Errorf("DataContext is empty")
	}
	document, _ := parseScope(datactx)
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

func parseFrontMatterArgs(args cty.Value, datactx plugin.MapData) (string, plugin.MapData, error) {
	format := args.GetAttr("format").AsString()

	if !slices.Contains(frontMatterAllowedFormats, format) {
		return "", nil, fmt.Errorf("invalid format: %s", format)
	}
	var data plugin.Data
	if datactx != nil {
		if qr, ok := datactx[frontMatterQueryResultKey]; ok {
			data = qr
		}
	}
	if data == nil {
		content := args.GetAttr("content")
		if !content.IsNull() {
			data = plugin.ConvertCtyToData(content)
		}
	}
	if data == nil {
		return "", nil, fmt.Errorf("%s and content are nil", frontMatterQueryResultKey)
	}
	m, ok := data.(plugin.MapData)
	if !ok {
		return "", nil, fmt.Errorf("invalid frontmatter data type: %T", data)
	}
	return format, m, nil
}

func renderFrontMatterContent(format string, m plugin.MapData) (string, error) {
	switch format {
	case "yaml":
		return renderYAMLFrontMatter(m)
	case "toml":
		return renderTOMLFrontMatter(m)
	case "json":
		return renderJSONFrontMatter(m)
	default:
		return "", fmt.Errorf("invalid format: %s", format)
	}
}

func renderYAMLFrontMatter(m plugin.MapData) (string, error) {
	var buf strings.Builder
	buf.WriteString("---\n")
	err := yaml.NewEncoder(&buf).Encode(m)
	if err != nil {
		return "", err
	}
	buf.WriteString("---")
	return buf.String(), nil
}

func renderTOMLFrontMatter(m plugin.MapData) (string, error) {
	var buf strings.Builder
	buf.WriteString("+++\n")
	err := toml.NewEncoder(&buf).Encode(m)
	if err != nil {
		return "", err
	}
	buf.WriteString("+++")
	return buf.String(), nil
}

func renderJSONFrontMatter(m plugin.MapData) (string, error) {
	var buf strings.Builder
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	err := enc.Encode(m)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
