package crowdstrike

import (
	"context"

	"github.com/crowdstrike/gofalcon/falcon"
	"github.com/crowdstrike/gofalcon/falcon/client/cspm_registration"
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/dataspec/constraint"
	"github.com/blackstork-io/fabric/plugin/plugindata"
)

func makeFalconCspmIomsDataSource(loader ClientLoaderFn) *plugin.DataSource {
	return &plugin.DataSource{
		Doc:      "The `falcon_cspm_ioms` data source fetches cloud indicators of misconfigurations (IOMs) from the Falcon security posture management (CSPM) feature",
		DataFunc: fetchFalconCspmIomsData(loader),
		Config:   makeDataSourceConfig(),
		Args: &dataspec.RootSpec{
			Attrs: []*dataspec.AttrSpec{
				{
					Name:        "size",
					Type:        cty.Number,
					Constraints: constraint.Integer | constraint.RequiredNonNull,
					Doc:         "limit the number of queried items",
				},
			},
		},
	}
}

func fetchFalconCspmIomsData(loader ClientLoaderFn) plugin.RetrieveDataFunc {
	return func(ctx context.Context, params *plugin.RetrieveDataParams) (plugindata.Data, diagnostics.Diag) {
		cli, err := loader(makeApiConfig(params.Config))
		if err != nil {
			return nil, diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Unable to create falcon client",
				Detail:   err.Error(),
			}}
		}
		size, _ := params.Args.GetAttrVal("size").AsBigFloat().Int64()
		apiParams := cspm_registration.NewGetConfigurationDetectionsParams().WithDefaults()
		apiParams.SetLimit(&size)
		apiParams.Context = ctx
		response, err := cli.CspmRegistration().GetConfigurationDetections(apiParams)
		if err != nil {
			return nil, diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Failed to fetch Falcon CSPM IOMs",
				Detail:   err.Error(),
			}}
		}
		if err = falcon.AssertNoError(response.GetPayload().Errors); err != nil {
			return nil, diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Failed to fetch Falcon CSPM IOMs",
				Detail:   err.Error(),
			}}
		}
		events := response.GetPayload().Resources.Events
		data, err := encodeResponse(events)
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
