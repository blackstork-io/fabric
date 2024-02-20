package builtin

import (
	"bytes"
	"context"
	"fmt"
	"slices"
	"strings"
	"text/template"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/plugin"
)

var textAllowedFormats = []string{"text", "title", "code", "blockquote"}

const (
	textMinAbsoluteTitleSize     = int64(1)
	textMaxAbsoluteTitleSize     = int64(6)
	textDefaultAbsoluteTitleSize = int64(1)
	textDefaultFormat            = "text"
	textDefaultCodeLanguage      = ""
)

func makeTextContentProvider() *plugin.ContentProvider {
	return &plugin.ContentProvider{
		ContentFunc: genTextContent,
		Args: hcldec.ObjectSpec{
			"text": &hcldec.AttrSpec{
				Name:     "text",
				Type:     cty.String,
				Required: true,
			},
			"format_as": &hcldec.AttrSpec{
				Name:     "format_as",
				Type:     cty.String,
				Required: false,
			},
			"absolute_title_size": &hcldec.AttrSpec{
				Name:     "absolute_title_size",
				Type:     cty.Number,
				Required: false,
			},
			"code_language": &hcldec.AttrSpec{
				Name:     "code_language",
				Type:     cty.String,
				Required: false,
			},
		},
	}
}

func genTextContent(ctx context.Context, params *plugin.ProvideContentParams) (*plugin.Content, hcl.Diagnostics) {
	text := params.Args.GetAttr("text")
	if text.IsNull() {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse arguments",
			Detail:   "text is required",
		}}
	}
	format := params.Args.GetAttr("format_as")
	if !format.IsNull() {
		if !slices.Contains(textAllowedFormats, format.AsString()) {
			return nil, hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Failed to parse arguments",
				Detail:   "format_as must be one of " + strings.Join(textAllowedFormats, ", "),
			}}
		}
	} else {
		format = cty.StringVal(textDefaultFormat)
	}
	absoluteTitleSize := params.Args.GetAttr("absolute_title_size")
	if absoluteTitleSize.IsNull() {
		absoluteTitleSize = cty.NumberIntVal(textDefaultAbsoluteTitleSize)
	}
	titleSize, _ := absoluteTitleSize.AsBigFloat().Int64()
	if titleSize < textMinAbsoluteTitleSize || titleSize > textMaxAbsoluteTitleSize {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse arguments",
			Detail:   fmt.Sprintf("absolute_title_size must be between %d and %d", textMinAbsoluteTitleSize, textMaxAbsoluteTitleSize),
		}}
	}
	codeLanguage := params.Args.GetAttr("code_language")
	if codeLanguage.IsNull() {
		codeLanguage = cty.StringVal(textDefaultCodeLanguage)
	}
	var (
		md  string
		err error
	)
	switch format.AsString() {
	case "text":
		md, err = genTextContentText(text.AsString(), params.DataContext)
	case "title":
		md, err = genTextContentTitle(text.AsString(), params.DataContext, titleSize)
	case "code":
		md, err = genTextContentCode(text.AsString(), params.DataContext, codeLanguage.AsString())
	case "blockquote":
		md, err = genTextContentBlockquote(text.AsString(), params.DataContext)
	default:
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse arguments",
			Detail:   fmt.Sprintf("format_as must be one of %s", strings.Join(textAllowedFormats, ", ")),
		}}
	}
	if err != nil {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to render text",
			Detail:   err.Error(),
		}}
	}
	return &plugin.Content{
		Markdown: md,
	}, nil
}

func genTextContentText(text string, datactx plugin.MapData) (string, error) {
	tmpl, err := template.New("text").Parse(text)
	if err != nil {
		return "", fmt.Errorf("failed to parse text template: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, datactx)
	if err != nil {
		return "", fmt.Errorf("failed to execute text template: %w", err)
	}
	return strings.TrimSpace(buf.String()), nil
}

func genTextContentTitle(text string, datactx plugin.MapData, titleSize int64) (string, error) {
	text, err := genTextContentText(text, datactx)
	if err != nil {
		return "", err
	}
	// remove all newlines
	text = strings.ReplaceAll(text, "\n", " ")
	return strings.Repeat("#", int(titleSize)) + " " + text, nil
}

func genTextContentCode(text string, datactx plugin.MapData, language string) (string, error) {
	text, err := genTextContentText(text, datactx)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("```%s\n%s\n```", language, text), nil
}

func genTextContentBlockquote(text string, datactx plugin.MapData) (string, error) {
	text, err := genTextContentText(text, datactx)
	if err != nil {
		return "", err
	}
	return "> " + strings.ReplaceAll(text, "\n", "\n> "), nil
}
