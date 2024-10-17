package crowdstrike

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/crowdstrike/gofalcon/falcon"
	"github.com/crowdstrike/gofalcon/falcon/client"
	"github.com/crowdstrike/gofalcon/falcon/client/cspm_registration"
	"github.com/crowdstrike/gofalcon/falcon/client/discover"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/dataspec/constraint"
	"github.com/blackstork-io/fabric/plugin/plugindata"
)

type CspmRegistrationClient interface {
	GetConfigurationDetections(params *cspm_registration.GetConfigurationDetectionsParams, opts ...cspm_registration.ClientOption) (*cspm_registration.GetConfigurationDetectionsOK, error)
}

type DiscoverClient interface {
	QueryHosts(params *discover.QueryHostsParams, opts ...discover.ClientOption) (*discover.QueryHostsOK, error)
	GetHosts(params *discover.GetHostsParams, opts ...discover.ClientOption) (*discover.GetHostsOK, error)
}

type Client interface {
	CspmRegistration() CspmRegistrationClient
	Discover() DiscoverClient
}

type ClientAdapter struct {
	client *client.CrowdStrikeAPISpecification
}

func (c *ClientAdapter) CspmRegistration() CspmRegistrationClient {
	return c.client.CspmRegistration
}

func (c *ClientAdapter) Discover() DiscoverClient {
	return c.client.Discover
}

type ClientLoaderFn func(cfg *falcon.ApiConfig) (client Client, err error)

var DefaultClientLoader = func(cfg *falcon.ApiConfig) (Client, error) {
	client, err := falcon.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &ClientAdapter{client}, nil
}

func Plugin(version string, loader ClientLoaderFn) *plugin.Schema {
	if loader == nil {
		loader = DefaultClientLoader
	}
	return &plugin.Schema{
		Name:    "blackstork/crowdstrike",
		Version: version,
		DataSources: plugin.DataSources{
			"falcon_cspm_ioms":             makeFalconCspmIomsDataSource(loader),
			"falcon_discover_host_details": makeFalconDiscoverHostDetailsDataSource(loader),
		},
	}
}

// shared config for all data sources
func makeDataSourceConfig() *dataspec.RootSpec {
	return &dataspec.RootSpec{
		Attrs: []*dataspec.AttrSpec{
			{
				Name:        "client_id",
				Type:        cty.String,
				Constraints: constraint.RequiredMeaningful,
				Doc:         "Client ID for accessing CrowdStrike Falcon Platform",
			},
			{
				Name:        "client_secret",
				Type:        cty.String,
				Constraints: constraint.RequiredMeaningful,
				Secret:      true,
				Doc:         "Client Secret for accessing CrowdStrike Falcon Platform",
			},
			{
				Name: "member_cid",
				Type: cty.String,
				Doc:  "Member CID for MSSP",
			},
			{
				Name: "client_cloud",
				Type: cty.String,
				OneOf: []cty.Value{
					cty.StringVal("autodiscover"),
					cty.StringVal("us-1"),
					cty.StringVal("us-2"),
					cty.StringVal("eu-1"),
					cty.StringVal("us-gov-1"),
					cty.StringVal("gov1"),
				},
				Doc:        "Falcon cloud abbreviation",
				ExampleVal: cty.StringVal("us-1"),
			},
		},
	}
}

func makeApiConfig(ctx context.Context, cfg *dataspec.Block) *falcon.ApiConfig {
	clientId := cfg.GetAttrVal("client_id").AsString()
	clientSecret := cfg.GetAttrVal("client_secret").AsString()
	apiCfg := &falcon.ApiConfig{
		ClientId:     clientId,
		ClientSecret: clientSecret,
	}
	memberCID := cfg.GetAttrVal("member_cid")
	if !memberCID.IsNull() {
		apiCfg.MemberCID = memberCID.AsString()
	}
	clientCloud := cfg.GetAttrVal("client_cloud")
	if !clientCloud.IsNull() {
		apiCfg.Cloud = falcon.Cloud(clientCloud.AsString())
	}
	apiCfg.Context = ctx
	return apiCfg
}

func encodeResponse(data any) (plugindata.Data, error) {
	raw, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to encode issue: %w", err)
	}
	return plugindata.UnmarshalJSON(raw)
}
