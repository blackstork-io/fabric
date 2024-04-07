package elastic

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/internal/elastic/kbclient"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
)

const (
	minSize     = 1
	defaultSize = 10
)

func makeElasticSecurityCasesDataSource(loader KibanaClientLoaderFn) *plugin.DataSource {
	return &plugin.DataSource{
		DataFunc: fetchElasticSecurityCases(loader),
		Config: dataspec.ObjectSpec{
			&dataspec.AttrSpec{
Name: "kibana_endpoint_url",
				Type:     cty.String,
				Required: true,
			},
			&dataspec.AttrSpec{
Name: "api_key_str",
				Type:     cty.String,
				Required: false,
			},
			&dataspec.AttrSpec{
Name: "api_key",
				Type:     cty.List(cty.String),
				Required: false,
			},
		},
		Args: dataspec.ObjectSpec{
			&dataspec.AttrSpec{
Name: "space_id",
				Type:     cty.String,
				Required: false,
			},
			&dataspec.AttrSpec{
Name: "assignees",
				Type:     cty.List(cty.String),
				Required: false,
			},
			&dataspec.AttrSpec{
Name: "default_search_operator",
				Type:     cty.String,
				Required: false,
			},
			&dataspec.AttrSpec{
Name: "from",
				Type:     cty.String,
				Required: false,
			},
			&dataspec.AttrSpec{
Name: "owner",
				Type:     cty.List(cty.String),
				Required: false,
			},
			&dataspec.AttrSpec{
Name: "reporters",
				Type:     cty.List(cty.String),
				Required: false,
			},
			&dataspec.AttrSpec{
Name: "search",
				Type:     cty.String,
				Required: false,
			},
			&dataspec.AttrSpec{
Name: "search_fields",
				Type:     cty.List(cty.String),
				Required: false,
			},
			&dataspec.AttrSpec{
Name: "severity",
				Type:     cty.String,
				Required: false,
			},
			&dataspec.AttrSpec{
Name: "sort_field",
				Type:     cty.String,
				Required: false,
			},
			&dataspec.AttrSpec{
Name: "sort_order",
				Type:     cty.String,
				Required: false,
			},
			&dataspec.AttrSpec{
Name: "status",
				Type:     cty.String,
				Required: false,
			},
			&dataspec.AttrSpec{
Name: "tags",
				Type:     cty.List(cty.String),
				Required: false,
			},
			&dataspec.AttrSpec{
Name: "to",
				Type:     cty.String,
				Required: false,
			},
			&dataspec.AttrSpec{
Name: "size",
				Type:     cty.Number,
				Required: false,
			},
		},
	}
}

func fetchElasticSecurityCases(loader KibanaClientLoaderFn) plugin.RetrieveDataFunc {
	return func(ctx context.Context, params *plugin.RetrieveDataParams) (plugin.Data, hcl.Diagnostics) {
		client, err := parseSecurityCasesConfig(loader, params.Config)
		if err != nil {
			return nil, hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Failed to parse configuration",
				Detail:   err.Error(),
			}}
		}
		req, err := parseSecurityCasesArgs(params.Args)
		if err != nil {
			return nil, hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Failed to parse arguments",
				Detail:   err.Error(),
			}}
		}
		size := defaultSize
		if attr := params.Args.GetAttr("size"); !attr.IsNull() {
			num, _ := attr.AsBigFloat().Int64()
			size = int(num)
			if size < minSize {
				return nil, hcl.Diagnostics{{
					Severity: hcl.DiagError,
					Summary:  "Invalid size",
					Detail:   "size must be greater than 0",
				}}
			}
		}
		req.PerPage = size
		req.Page = 1
		cases := plugin.ListData{}
		for {
			res, err := client.ListSecurityCases(ctx, req)
			if err != nil {
				return nil, hcl.Diagnostics{{
					Severity: hcl.DiagError,
					Summary:  "Failed to list security cases",
					Detail:   err.Error(),
				}}
			}
			for _, c := range res.Cases {
				data, err := plugin.ParseDataAny(c)
				if err != nil {
					return nil, hcl.Diagnostics{{
						Severity: hcl.DiagError,
						Summary:  "Failed to parse security case",
						Detail:   err.Error(),
					}}
				}
				cases = append(cases, data)
			}
			if len(cases) >= size || req.Page*req.PerPage >= res.Total {
				break
			}
			req.Page++
		}
		return cases, nil
	}
}

func parseSecurityCasesArgs(args cty.Value) (*kbclient.ListSecurityCasesReq, error) {
	if args.IsNull() {
		return nil, fmt.Errorf("arguments are required")
	}
	req := &kbclient.ListSecurityCasesReq{}
	if attr := args.GetAttr("space_id"); !attr.IsNull() {
		req.SpaceID = kbclient.String(attr.AsString())
	}
	if attr := args.GetAttr("assignees"); !attr.IsNull() {
		list := []string{}
		for _, v := range attr.AsValueSlice() {
			list = append(list, v.AsString())
		}
		req.Assignees = list
	}
	if attr := args.GetAttr("default_search_operator"); !attr.IsNull() {
		req.DefaultSearchOperator = kbclient.String(attr.AsString())
	}
	if attr := args.GetAttr("from"); !attr.IsNull() {
		req.From = kbclient.String(attr.AsString())
	}
	if attr := args.GetAttr("owner"); !attr.IsNull() {
		list := []string{}
		for _, v := range attr.AsValueSlice() {
			list = append(list, v.AsString())
		}
		req.Owner = list
	}
	if attr := args.GetAttr("reporters"); !attr.IsNull() {
		list := []string{}
		for _, v := range attr.AsValueSlice() {
			list = append(list, v.AsString())
		}
		req.Reporters = list
	}
	if attr := args.GetAttr("search"); !attr.IsNull() {
		req.Search = kbclient.String(attr.AsString())
	}
	if attr := args.GetAttr("search_fields"); !attr.IsNull() {
		list := []string{}
		for _, v := range attr.AsValueSlice() {
			list = append(list, v.AsString())
		}
		req.SearchFields = list
	}
	if attr := args.GetAttr("severity"); !attr.IsNull() {
		req.Severity = kbclient.String(attr.AsString())
	}
	if attr := args.GetAttr("sort_field"); !attr.IsNull() {
		req.SortField = kbclient.String(attr.AsString())
	}
	if attr := args.GetAttr("sort_order"); !attr.IsNull() {
		req.SortOrder = kbclient.String(attr.AsString())
	}
	if attr := args.GetAttr("status"); !attr.IsNull() {
		req.Status = kbclient.String(attr.AsString())
	}
	if attr := args.GetAttr("tags"); !attr.IsNull() {
		list := []string{}
		for _, v := range attr.AsValueSlice() {
			list = append(list, v.AsString())
		}
		req.Tags = list
	}
	if attr := args.GetAttr("to"); !attr.IsNull() {
		req.To = kbclient.String(attr.AsString())
	}
	return req, nil
}

func parseSecurityCasesConfig(loader KibanaClientLoaderFn, cfg cty.Value) (kbclient.Client, error) {
	if cfg.IsNull() {
		return nil, fmt.Errorf("configuration is required")
	}
	var url string
	var apiKey *string
	if attr := cfg.GetAttr("kibana_endpoint_url"); attr.IsNull() {
		return nil, fmt.Errorf("kibana_endpoint_url is required")
	} else {
		url = attr.AsString()
	}
	if attr := cfg.GetAttr("api_key_str"); !attr.IsNull() {
		apiKey = kbclient.String(attr.AsString())
	} else {
		if attr := cfg.GetAttr("api_key"); !attr.IsNull() {
			list := attr.AsValueSlice()
			if len(list) != 2 {
				return nil, fmt.Errorf("api_key must be a list of 2 strings")
			}
			key := base64.RawURLEncoding.EncodeToString([]byte(list[0].AsString() + ":" + list[1].AsString()))
			apiKey = kbclient.String(key)
		}
	}
	return loader(url, apiKey), nil
}
