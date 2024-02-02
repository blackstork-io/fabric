package json

import (
	"os"

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
			Kind:       "data",
			Name:       "json",
			Version:    plugininterface.Version(*Version),
			ConfigSpec: nil,
			InvocationSpec: &hcldec.ObjectSpec{
				"glob": &hcldec.AttrSpec{
					Name:     "glob",
					Type:     cty.String,
					Required: true,
				},
			},
		},
	}
}

func (Plugin) Call(args plugininterface.Args) plugininterface.Result {
	glob := args.Args.GetAttr("glob").AsString()
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
	docs, err := readFS(filesystem, glob)
	if err != nil {
		return plugininterface.Result{
			Diags: hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Failed to read json files",
				Detail:   err.Error(),
			}},
		}
	}
	data := make([]any, len(docs))
	for i, doc := range docs {
		data[i] = doc.Map()
	}
	return plugininterface.Result{
		Result: data,
	}
}
