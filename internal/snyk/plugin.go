package snyk

import (
	"fmt"

	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/internal/snyk/client"
	"github.com/blackstork-io/fabric/plugin"
)

const (
	pageSize = 100
)

type ClientLoadFn func(apiKey string) client.Client

var DefaultClientLoader ClientLoadFn = client.New

func Plugin(version string, loader ClientLoadFn) *plugin.Schema {
	if loader == nil {
		loader = DefaultClientLoader
	}
	return &plugin.Schema{
		Name:    "blackstork/snyk",
		Version: version,
		DataSources: plugin.DataSources{
			"snyk_issues": makeSnykIssuesDataSource(loader),
		},
	}
}

func makeClient(loader ClientLoadFn, cfg cty.Value) (client.Client, error) {
	if cfg.IsNull() {
		return nil, fmt.Errorf("configuration is required")
	}
	apiKey := cfg.GetAttr("api_key")
	if apiKey.IsNull() || apiKey.AsString() == "" {
		return nil, fmt.Errorf("api_key is required in configuration")
	}
	return loader(apiKey.AsString()), nil
}
