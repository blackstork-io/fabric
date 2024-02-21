package splunk

import (
	"fmt"

	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/internal/splunk/client"
	"github.com/blackstork-io/fabric/plugin"
)

type ClientLoadFn func(token, host, deployment string) client.Client

var DefaultClientLoader ClientLoadFn = client.New

func Plugin(version string, loader ClientLoadFn) *plugin.Schema {
	if loader == nil {
		loader = DefaultClientLoader
	}
	return &plugin.Schema{
		Name:    "blackstork/splunk",
		Version: version,
		DataSources: plugin.DataSources{
			"splunk_search": makeSplunkSearchDataSchema(loader),
		},
	}
}

func makeClient(loader ClientLoadFn, cfg cty.Value) (client.Client, error) {
	if cfg.IsNull() {
		return nil, fmt.Errorf("configuration is required")
	}

	token := cfg.GetAttr("auth_token")
	if token.IsNull() || token.AsString() == "" {
		return nil, fmt.Errorf("auth_token is required in configuration")
	}
	host := cfg.GetAttr("host")
	if host.IsNull() {
		host = cty.StringVal("")
	}
	deployment := cfg.GetAttr("deployment_name")
	if deployment.IsNull() {
		deployment = cty.StringVal("")
	}
	if host.AsString() == "" && deployment.AsString() == "" {
		return nil, fmt.Errorf("host or deployment_name is required in configuration")
	}
	return loader(token.AsString(), host.AsString(), deployment.AsString()), nil
}
