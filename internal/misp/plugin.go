package misp

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/internal/misp/client"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/dataspec/constraint"
	"github.com/blackstork-io/fabric/plugin/plugindata"
)

type Client interface {
	RestSearchEvents(ctx context.Context, req client.RestSearchEventsRequest) (events client.RestSearchEventsResponse, err error)
}

type ClientLoaderFn func(cfg *dataspec.Block) Client

func DefaultClientLoader(cfg *dataspec.Block) Client {
	apiKey := cfg.GetAttrVal("api_key").AsString()
	baseUrl := cfg.GetAttrVal("base_url").AsString()
	skipSsl := cfg.GetAttrVal("skip_ssl").True()
	opts := []client.ClientOption{}
	if skipSsl {
		cli := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
		opts = append(opts, client.WithHTTPClient(cli))
	}
	return client.NewClient(baseUrl, apiKey, opts...)
}

func Plugin(version string, loader ClientLoaderFn) *plugin.Schema {
	if loader == nil {
		loader = DefaultClientLoader
	}
	return &plugin.Schema{
		Name:    "blackstork/misp",
		Version: version,
		DataSources: plugin.DataSources{
			"misp_events": makeMispEventsDataSource(loader),
		},
	}
}

// shared config for all data sources
func makeDataSourceConfig() *dataspec.RootSpec {
	return &dataspec.RootSpec{
		Attrs: []*dataspec.AttrSpec{
			{
				Name:        "api_key",
				Type:        cty.String,
				Constraints: constraint.RequiredMeaningful,
				Doc:         "misp api key",
			},
			{
				Name:        "base_url",
				Type:        cty.String,
				Constraints: constraint.RequiredMeaningful,
				Doc:         "misp base url",
			},
			{
				Name:       "skip_ssl",
				Type:       cty.Bool,
				Doc:        "skip ssl verification",
				DefaultVal: cty.BoolVal(false),
			},
		},
	}
}

func encodeResponse(data any) (plugindata.Data, error) {
	raw, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to encode issue: %w", err)
	}
	return plugindata.UnmarshalJSON(raw)
}
