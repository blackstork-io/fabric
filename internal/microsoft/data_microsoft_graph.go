package microsoft

import (
	"context"
	"net/url"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/dataspec/constraint"
	"github.com/blackstork-io/fabric/plugin/plugindata"
)

func makeMicrosoftGraphDataSource(loader MicrosoftGraphClientLoadFn) *plugin.DataSource {
	return &plugin.DataSource{
		Doc:      "The `microsoft_graph` data source queries Microsoft Graph.",
		DataFunc: fetchMicrosoftGraph(loader),
		Config: &dataspec.RootSpec{
			Attrs: []*dataspec.AttrSpec{
				{
					Doc:         "The Azure client ID",
					Name:        "client_id",
					Type:        cty.String,
					Constraints: constraint.RequiredNonNull,
				},
				{
					Doc:    "The Azure client secret. Required if private_key_file/privat_key/cert_thumbprint is not provided.",
					Name:   "client_secret",
					Type:   cty.String,
					Secret: true,
				},
				{
					Doc:         "The Azure tenant ID",
					Name:        "tenant_id",
					Type:        cty.String,
					Constraints: constraint.RequiredNonNull,
				},
				{
					Doc:  "The path to the private key file. Ignored if private_key/client_secret is provided.",
					Name: "private_key_file",
					Type: cty.String,
				},
				{
					Doc:  "The private key contents. Ignored if client_secret is provided.",
					Name: "private_key",
					Type: cty.String,
				},
				{
					Doc:  "The key passphrase. Ignored if client_secret is provided.",
					Name: "key_passphrase",
					Type: cty.String,
				},
			},
		},
		Args: &dataspec.RootSpec{
			Attrs: []*dataspec.AttrSpec{
				{
					Doc:        "The API version",
					Name:       "api_version",
					Type:       cty.String,
					DefaultVal: cty.StringVal("beta"),
				},
				{
					Doc:         "The endpoint to query",
					Name:        "endpoint",
					Type:        cty.String,
					Constraints: constraint.RequiredNonNull,
					ExampleVal:  cty.StringVal("/security/incidents"),
				},
				{
					Doc:  "The query parameters",
					Name: "query_params",
					Type: cty.Map(cty.String),
				},
			},
		},
	}
}

func fetchMicrosoftGraph(loader MicrosoftGraphClientLoadFn) plugin.RetrieveDataFunc {
	return func(ctx context.Context, params *plugin.RetrieveDataParams) (plugindata.Data, diagnostics.Diag) {
		apiVersion := params.Args.GetAttrVal("api_version").AsString()
		cli, err := loader(ctx, apiVersion, params.Config)
		if err != nil {
			return nil, diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Unable to create microsoft graph client",
				Detail:   err.Error(),
			}}
		}
		endPoint := params.Args.GetAttrVal("endpoint").AsString()
		queryParamsAttr := params.Args.GetAttrVal("query_params")
		var queryParams url.Values

		if !queryParamsAttr.IsNull() {
			queryParams = url.Values{}
			queryMap := queryParamsAttr.AsValueMap()
			for k, v := range queryMap {
				queryParams.Add(k, v.AsString())
			}
		}

		response, err := cli.QueryGraph(ctx, endPoint, queryParams)
		if err != nil {
			return nil, diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Failed to query microsoft graph",
				Detail:   err.Error(),
			}}
		}
		data, err := plugindata.ParseAny(response)
		if err != nil {
			return nil, diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Failed to parse response",
				Detail:   err.Error(),
			}}
		}
		return data, nil
	}
}
