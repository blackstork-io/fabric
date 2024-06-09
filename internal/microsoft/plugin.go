package microsoft

import (
	"context"
	"fmt"

	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/internal/microsoft/client"
	"github.com/blackstork-io/fabric/plugin"
)

type ClientLoadFn func() client.Client

var DefaultClientLoader ClientLoadFn = client.New

func Plugin(version string, loader ClientLoadFn) *plugin.Schema {
	if loader == nil {
		loader = DefaultClientLoader
	}
	return &plugin.Schema{
		Doc:     "The `microsoft` plugin for Microsoft services.",
		Name:    "blackstork/microsoft",
		Version: version,
		DataSources: plugin.DataSources{
			"microsoft_sentinel_incidents": makeMicrosoftSentinelIncidentsDataSource(loader),
		},
	}
}

func makeClient(ctx context.Context, loader ClientLoadFn, cfg cty.Value) (client.Client, error) {
	if cfg.IsNull() {
		return nil, fmt.Errorf("configuration is required")
	}
	cli := loader()
	res, err := cli.GetClientCredentialsToken(ctx, &client.GetClientCredentialsTokenReq{
		TenantID:     cfg.GetAttr("tenant_id").AsString(),
		ClientID:     cfg.GetAttr("client_id").AsString(),
		ClientSecret: cfg.GetAttr("client_secret").AsString(),
	})
	if err != nil {
		return nil, err
	}
	cli.UseAuth(res.AccessToken)
	return cli, nil
}
