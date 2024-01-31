package shout

import (
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	plugininterface "github.com/blackstork-io/fabric/plugininterface/v1"
)

var Version = semver.MustParse("0.1.2")

type Plugin struct{}

func (Plugin) GetPlugins() []plugininterface.Plugin {
	return []plugininterface.Plugin{
		{
			Namespace:  "blackstork-example",
			Kind:       "content",
			Name:       "shout",
			Version:    plugininterface.Version(*Version),
			ConfigSpec: nil,
			InvocationSpec: &hcldec.ObjectSpec{
				"text": &hcldec.AttrSpec{
					Name:     "text",
					Type:     cty.String,
					Required: true,
				},
			},
		},
	}
}

func (p Plugin) Call(args plugininterface.Args) plugininterface.Result {
	text := args.Args.GetAttr("text")

	return plugininterface.Result{
		Result: strings.ToUpper(text.AsString()),
	}
}
