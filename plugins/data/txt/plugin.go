package txt

import (
	"io/fs"
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
			Name:       "txt",
			Version:    plugininterface.Version(*Version),
			ConfigSpec: nil,
			InvocationSpec: &hcldec.ObjectSpec{
				"path": &hcldec.AttrSpec{
					Name:     "path",
					Type:     cty.String,
					Required: true,
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
	f, err := filesystem.Open(path.AsString())
	if err != nil {
		return plugininterface.Result{
			Diags: hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Failed to open txt file",
				Detail:   err.Error(),
			}},
		}
	}
	defer f.Close()
	data, err := fs.ReadFile(filesystem, path.AsString())
	if err != nil {
		return plugininterface.Result{
			Diags: hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Failed to read txt file",
				Detail:   err.Error(),
			}},
		}
	}
	return plugininterface.Result{
		Result: string(data),
	}
}
