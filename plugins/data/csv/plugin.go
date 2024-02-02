package csv

import (
	"os"

	"github.com/Masterminds/semver/v3"
	"github.com/blackstork-io/fabric/plugininterface/v1"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"
)

var Version = semver.MustParse("0.1.0")

const defaultDelimiter = ','

type Plugin struct{}

func (Plugin) GetPlugins() []plugininterface.Plugin {
	return []plugininterface.Plugin{
		{
			Namespace:  "blackstork",
			Kind:       "data",
			Name:       "csv",
			Version:    plugininterface.Version(*Version),
			ConfigSpec: nil,
			InvocationSpec: &hcldec.ObjectSpec{
				"path": &hcldec.AttrSpec{
					Name:     "path",
					Type:     cty.String,
					Required: true,
				},
				"delimiter": &hcldec.AttrSpec{
					Name:     "delimiter",
					Type:     cty.String,
					Required: false,
				},
			},
		},
	}
}

func (Plugin) Call(args plugininterface.Args) plugininterface.Result {
	path := args.Args.GetAttr("path")
	if path.IsNull() || path.AsString() == "" {
		return plugininterface.Result{
			Diags: hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "path is required",
			}},
		}
	}
	delim := args.Args.GetAttr("delimiter")
	if delim.IsNull() {
		delim = cty.StringVal(string(defaultDelimiter))
	}
	if len(delim.AsString()) != 1 {
		return plugininterface.Result{
			Diags: hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "delimiter must be a single character",
			}},
		}
	}
	delimRune := []rune(delim.AsString())[0]
	wd, err := os.Getwd()
	if err != nil {
		return plugininterface.Result{
			Diags: hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Failed to get current working directory",
				Detail:   err.Error(),
			}},
		}
	}
	filesystem := os.DirFS(wd)
	data, err := readFS(filesystem, path.AsString(), delimRune)
	if err != nil {
		return plugininterface.Result{
			Diags: hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Failed to read csv file",
				Detail:   err.Error(),
			}},
		}
	}
	return plugininterface.Result{
		Result: data,
	}
}
