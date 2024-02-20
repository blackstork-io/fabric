package elasticsearch

import (
	"context"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/plugin"
)

func makeElasticSearchDataSource() *plugin.DataSource {
	return &plugin.DataSource{
		DataFunc: fetchElasticSearchData,
		Config: hcldec.ObjectSpec{
			"base_url": &hcldec.AttrSpec{
				Name:     "base_url",
				Type:     cty.String,
				Required: false,
			},
			"cloud_id": &hcldec.AttrSpec{
				Name:     "cloud_id",
				Type:     cty.String,
				Required: false,
			},
			"api_key_str": &hcldec.AttrSpec{
				Name:     "api_key_str",
				Type:     cty.String,
				Required: false,
			},
			"api_key": &hcldec.AttrSpec{
				Name:     "api_key",
				Type:     cty.List(cty.String),
				Required: false,
			},
			"basic_auth_username": &hcldec.AttrSpec{
				Name:     "basic_auth_username",
				Type:     cty.String,
				Required: false,
			},
			"basic_auth_password": &hcldec.AttrSpec{
				Name:     "basic_auth_password",
				Type:     cty.String,
				Required: false,
			},
			"bearer_auth": &hcldec.AttrSpec{
				Name:     "bearer_auth",
				Type:     cty.String,
				Required: false,
			},
			"ca_certs": &hcldec.AttrSpec{
				Name:     "ca_certs",
				Type:     cty.String,
				Required: false,
			},
		},
		Args: hcldec.ObjectSpec{
			"index": &hcldec.AttrSpec{
				Name:     "index",
				Type:     cty.String,
				Required: true,
			},
			"id": &hcldec.AttrSpec{
				Name:     "id",
				Type:     cty.String,
				Required: false,
			},
			"query_string": &hcldec.AttrSpec{
				Name:     "query_string",
				Type:     cty.String,
				Required: false,
			},
			"query": &hcldec.AttrSpec{
				Name:     "query",
				Type:     cty.Map(cty.DynamicPseudoType),
				Required: false,
			},
			"fields": &hcldec.AttrSpec{
				Name:     "fields",
				Type:     cty.List(cty.String),
				Required: false,
			},
			"size": &hcldec.AttrSpec{
				Name:     "size",
				Type:     cty.Number,
				Required: false,
			},
		},
	}
}

func fetchElasticSearchData(ctx context.Context, params *plugin.RetrieveDataParams) (plugin.Data, hcl.Diagnostics) {
	client, err := makeClient(params.Config)
	if err != nil {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to create elasticsearch client",
			Detail:   err.Error(),
		}}
	}
	id := params.Args.GetAttr("id")
	var data plugin.Data
	if !id.IsNull() {
		data, err = getByID(client.Get, params.Args)
	} else {
		data, err = search(client.Search, params.Args)
	}
	if err != nil {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to get data",
			Detail:   err.Error(),
		}}
	}
	return data, nil
}
