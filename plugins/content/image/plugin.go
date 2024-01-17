package image

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/blackstork-io/fabric/plugininterface/v1"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"
)

var Version = semver.MustParse("0.1.0")

type Plugin struct{}

func (Plugin) GetPlugins() []plugininterface.Plugin {
	return []plugininterface.Plugin{
		{
			Namespace:  "blackstork",
			Kind:       "content",
			Name:       "image",
			Version:    plugininterface.Version(*Version),
			ConfigSpec: nil,
			InvocationSpec: &hcldec.ObjectSpec{
				"src": &hcldec.AttrSpec{
					Name:     "src",
					Type:     cty.String,
					Required: true,
				},
				"alt": &hcldec.AttrSpec{
					Name:     "alt",
					Type:     cty.String,
					Required: false,
				},
			},
		},
	}
}

func (p Plugin) Call(args plugininterface.Args) plugininterface.Result {
	src, alt, err := p.parseArgs(args)
	if err != nil {
		return plugininterface.Result{
			Diags: hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Failed to parse arguments",
				Detail:   err.Error(),
			}},
		}
	}
	return plugininterface.Result{
		Result: p.render(src, alt),
	}
}

func (p Plugin) parseArgs(args plugininterface.Args) (string, string, error) {
	src := args.Args.GetAttr("src")
	if src.IsNull() || src.AsString() == "" {
		return "", "", errors.New("src is required")
	}
	alt := args.Args.GetAttr("alt")
	if alt.IsNull() {
		alt = cty.StringVal("")
	}
	return src.AsString(), alt.AsString(), nil
}

// render markdown image
func (p Plugin) render(src, alt string) string {
	src = strings.TrimSpace(strings.ReplaceAll(src, "\n", ""))
	alt = strings.TrimSpace(strings.ReplaceAll(alt, "\n", ""))
	return fmt.Sprintf("![%s](%s)", alt, src)
}
