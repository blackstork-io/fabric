package builtin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/pelletier/go-toml/v2"
	"github.com/zclconf/go-cty/cty"
	"gopkg.in/yaml.v3"

	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
)

const (
	frontMatterDefaultFormat  = "yaml"
	frontMatterQueryResultKey = "query_result"
)

var frontMatterAllowedFormats = []string{"yaml", "toml", "json"}

func makeFrontMatterContentProvider() *plugin.ContentProvider {
	return &plugin.ContentProvider{
		ContentFunc: genFrontMatterContent,
		Args: dataspec.ObjectSpec{
			&dataspec.AttrSpec{
				Name:     "format",
				Type:     cty.String,
				Required: false,
			},
			&dataspec.AttrSpec{
				Name:     "content",
				Type:     cty.Map(cty.DynamicPseudoType),
				Required: false,
			},
		},
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
	format := args.GetAttr("format")
	if format.IsNull() || format.AsString() == "" {
		format = cty.StringVal(frontMatterDefaultFormat)
	}
	if !slices.Contains(frontMatterAllowedFormats, format.AsString()) {
		return "", nil, fmt.Errorf("invalid format: %s", format.AsString())
	}
	var m plugin.MapData
	if datactx != nil {
		if queryResult, ok := datactx[frontMatterQueryResultKey]; ok {
			if qr, ok := queryResult.(plugin.MapData); ok {
				m = qr
			} else {
				return "", nil, fmt.Errorf("invalid query result: %T", queryResult)
			}
		}
	}
	content := args.GetAttr("content")
	if m == nil {
		if !content.IsNull() {
			m = convertCtyToDataMap(content)
		} else {
			return "", nil, errors.New("query_result and content are nil")
		}
	}
	return format.AsString(), m, nil
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
	buf.WriteString("---\n")
	return buf.String(), nil
}

func renderTOMLFrontMatter(m plugin.MapData) (string, error) {
	var buf strings.Builder
	buf.WriteString("+++\n")
	err := toml.NewEncoder(&buf).Encode(m)
	if err != nil {
		return "", err
	}
	buf.WriteString("+++\n")
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

func convertCtyToDataMap(v cty.Value) plugin.MapData {
	result := make(plugin.MapData)
	for k, v := range v.AsValueMap() {
		result[k] = convertCtyToData(v)
	}
	return result
}

func convertCtyToData(v cty.Value) plugin.Data {
	if v.IsNull() {
		return nil
	}
	t := v.Type()
	switch {
	case t == cty.String:
		return plugin.StringData(v.AsString())
	case t == cty.Number:
		if v.AsBigFloat().IsInt() {
			n, _ := v.AsBigFloat().Float64()
			return plugin.NumberData(n)
		}
	case t == cty.Bool:
		return plugin.BoolData(v.True())
	case t.IsMapType() || t.IsObjectType():
		return convertCtyToDataMap(v)
	case t.IsListType():
		return convertCtyToDataList(v)
	}
	return nil
}

func convertCtyToDataList(v cty.Value) plugin.ListData {
	var result plugin.ListData
	for _, v := range v.AsValueSlice() {
		result = append(result, convertCtyToData(v))
	}
	return result
}
