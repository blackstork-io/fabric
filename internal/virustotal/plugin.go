package virustotal

import (
	"fmt"

	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/internal/virustotal/client"
	"github.com/blackstork-io/fabric/plugin"
)

type ClientLoadFn func(key string) client.Client

var DefaultClientLoader ClientLoadFn = client.New

func Plugin(version string, loader ClientLoadFn) *plugin.Schema {
	if loader == nil {
		loader = DefaultClientLoader
	}
	return &plugin.Schema{
		Name:    "blackstork/virustotal",
		Version: version,
		DataSources: plugin.DataSources{
			"virustotal_api_usage": makeVirusTotalAPIUsageDataSchema(loader),
		},
	}
}

func makeClient(loader ClientLoadFn, cfg cty.Value) (client.Client, error) {
	if cfg.IsNull() {
		return nil, fmt.Errorf("configuration is required")
	}
	key := cfg.GetAttr("api_key")
	if key.IsNull() || key.AsString() == "" {
		return nil, fmt.Errorf("api_key is required in configuration")
	}
	return loader(key.AsString()), nil
}
