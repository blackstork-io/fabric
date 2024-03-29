package opencti

import (
	"context"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/plugin"
)

func makeOpenCTIDataSource() *plugin.DataSource {
	return &plugin.DataSource{
		Config: hcldec.ObjectSpec{
			"graphql_url": &hcldec.AttrSpec{
				Name:     "graphql_url",
				Type:     cty.String,
				Required: true,
			},
			"auth_token": &hcldec.AttrSpec{
				Name:     "auth_token",
				Type:     cty.String,
				Required: false,
			},
		},
		Args: hcldec.ObjectSpec{
			"graphql_query": &hcldec.AttrSpec{
				Name:     "graphql_query",
				Type:     cty.String,
				Required: true,
			},
		},
		DataFunc: fetchOpenCTIData,
	}
}

func fetchOpenCTIData(ctx context.Context, params *plugin.RetrieveDataParams) (plugin.Data, hcl.Diagnostics) {
	url := params.Config.GetAttr("graphql_url")
	if url.IsNull() || url.AsString() == "" {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse config",
			Detail:   "graphql_url is required",
		}}
	}
	authToken := params.Config.GetAttr("auth_token")
	if authToken.IsNull() {
		authToken = cty.StringVal("")
	}
	query := params.Args.GetAttr("graphql_query")
	if query.IsNull() || query.AsString() == "" {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse arguments",
			Detail:   "graphql_query is required",
		}}
	}
	if err := validateQuery(query.AsString()); err != nil {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Invalid GraphQL query",
			Detail:   err.Error(),
		}}
	}
	result, err := executeQuery(ctx, url.AsString(), query.AsString(), authToken.AsString())
	if err != nil {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to execute query",
			Detail:   err.Error(),
		}}
	}
	return result, nil
}
