package sentinel

import (
	"fmt"

	"github.com/blackstork-io/fabric/internal/sentinel/client"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/zclconf/go-cty/cty"
)

type ClientLoadFn func(token string) client.Client

var DefaultClientLoader ClientLoadFn = client.New

func Plugin(version string, loader ClientLoadFn) *plugin.Schema {
	if loader == nil {
		loader = DefaultClientLoader
	}
	return &plugin.Schema{
		Name:    "blackstork/microsoft_sentinel",
		Version: version,
		DataSources: plugin.DataSources{
			"microsoft_sentinel_incidents": makeMicrosoftSentinelIncidentsDataSource(loader),
		},
	}
}

func makeClient(loader ClientLoadFn, cfg cty.Value) (client.Client, error) {
	if cfg.IsNull() {
		return nil, fmt.Errorf("configuration is required")
	}
	return loader(""), nil
}
