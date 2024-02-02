package text

import (
	"bytes"
	"errors"
	"fmt"
	"slices"
	"strings"
	"text/template"

	"github.com/Masterminds/semver/v3"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/plugininterface/v1"
)

var Version = semver.MustParse("0.1.0")
var allowedFormats = []string{"text", "title", "code", "blockquote"}

const (
	minAbsoluteTitleSize     = int64(1)
	maxAbsoluteTitleSize     = int64(6)
	defaultAbsoluteTitleSize = int64(1)
	defaultFormat            = "text"
	defaultCodeLanguage      = ""
)

type Plugin struct{}

func (Plugin) GetPlugins() []plugininterface.Plugin {
	return []plugininterface.Plugin{
		{
			Namespace:  "blackstork",
			Kind:       "content",
			Name:       "text",
			Version:    plugininterface.Version(*Version),
			ConfigSpec: nil,
			InvocationSpec: &hcldec.ObjectSpec{
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
		},
	}
}

func (p Plugin) render(args cty.Value, datactx map[string]any) (string, error) {
	text := args.GetAttr("text")
	if text.IsNull() {
		return "", errors.New("text is required")
	}
	format := args.GetAttr("format_as")
	if !format.IsNull() {
		if !slices.Contains(allowedFormats, format.AsString()) {
			return "", errors.New("format_as must be one of " + strings.Join(allowedFormats, ", "))
		}
	} else {
		format = cty.StringVal(defaultFormat)
	}
	absoluteTitleSize := args.GetAttr("absolute_title_size")
	if absoluteTitleSize.IsNull() {
		absoluteTitleSize = cty.NumberIntVal(defaultAbsoluteTitleSize)
	}
	titleSize, _ := absoluteTitleSize.AsBigFloat().Int64()
	if titleSize < minAbsoluteTitleSize || titleSize > maxAbsoluteTitleSize {
		return "", fmt.Errorf("absolute_title_size must be between %d and %d", minAbsoluteTitleSize, maxAbsoluteTitleSize)
	}
	codeLanguage := args.GetAttr("code_language")
	if codeLanguage.IsNull() {
		codeLanguage = cty.StringVal(defaultCodeLanguage)
	}
	switch format.AsString() {
	case "text":
		return p.renderText(text.AsString(), datactx)
	case "title":
		return p.renderTitle(text.AsString(), datactx, titleSize)
	case "code":
		return p.renderCode(text.AsString(), datactx, codeLanguage.AsString())
	case "blockquote":
		return p.renderBlockquote(text.AsString(), datactx)
	}
	panic("unreachable")
}

func (p Plugin) renderText(text string, datactx map[string]any) (string, error) {
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

func (p Plugin) renderTitle(text string, datactx map[string]any, titleSize int64) (string, error) {
	text, err := p.renderText(text, datactx)
	if err != nil {
		return "", err
	}
	// remove all newlines
	text = strings.ReplaceAll(text, "\n", " ")
	return strings.Repeat("#", int(titleSize)) + " " + text, nil
}

func (p Plugin) renderCode(text string, datactx map[string]any, language string) (string, error) {
	text, err := p.renderText(text, datactx)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("```%s\n%s\n```", language, text), nil
}

func (p Plugin) renderBlockquote(text string, datactx map[string]any) (string, error) {
	text, err := p.renderText(text, datactx)
	if err != nil {
		return "", err
	}
	return "> " + strings.ReplaceAll(text, "\n", "\n> "), nil
}

func (p Plugin) Call(args plugininterface.Args) plugininterface.Result {
	result, err := p.render(args.Args, args.Context)
	if err != nil {
		return plugininterface.Result{
			Diags: hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Failed to render text",
				Detail:   err.Error(),
			}},
		}
	}
	return plugininterface.Result{
		Result: result,
	}
}
