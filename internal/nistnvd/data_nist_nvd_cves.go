package nistnvd

import (
	"context"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/internal/nistnvd/client"
	"github.com/blackstork-io/fabric/plugin"
)

const (
	defaultLimit = 1000
	minLimit     = 1
	maxLimit     = 2000
)

func makeNistNvdCvesDataSource(loader ClientLoadFn) *plugin.DataSource {
	return &plugin.DataSource{
		DataFunc: fetchNistNvdCvesData(loader),
		Config: hcldec.ObjectSpec{
			"api_key": &hcldec.AttrSpec{
				Name:     "api_key",
				Type:     cty.String,
				Required: false,
			},
		},
		Args: hcldec.ObjectSpec{
			"last_mod_start_date": &hcldec.AttrSpec{
				Name:     "last_mod_start_date",
				Type:     cty.String,
				Required: false,
			},
			"last_mod_end_date": &hcldec.AttrSpec{
				Name:     "last_mod_end_date",
				Type:     cty.String,
				Required: false,
			},
			"pub_start_date": &hcldec.AttrSpec{
				Name:     "pub_start_date",
				Type:     cty.String,
				Required: false,
			},
			"pub_end_date": &hcldec.AttrSpec{
				Name:     "pub_end_date",
				Type:     cty.String,
				Required: false,
			},
			"cpe_name": &hcldec.AttrSpec{
				Name:     "cpe_name",
				Type:     cty.String,
				Required: false,
			},
			"cve_id": &hcldec.AttrSpec{
				Name:     "cve_id",
				Type:     cty.String,
				Required: false,
			},
			"cvss_v3_metrics": &hcldec.AttrSpec{
				Name:     "cvss_v3_metrics",
				Type:     cty.String,
				Required: false,
			},
			"cvss_v3_severity": &hcldec.AttrSpec{
				Name:     "cvss_v3_severity",
				Type:     cty.String,
				Required: false,
			},
			"cwe_id": &hcldec.AttrSpec{
				Name:     "cwe_id",
				Type:     cty.String,
				Required: false,
			},
			"keyword_search": &hcldec.AttrSpec{
				Name:     "keyword_search",
				Type:     cty.String,
				Required: false,
			},
			"virtual_match_string": &hcldec.AttrSpec{
				Name:     "virtual_match_string",
				Type:     cty.String,
				Required: false,
			},
			"source_identifier": &hcldec.AttrSpec{
				Name:     "source_identifier",
				Type:     cty.String,
				Required: false,
			},
			"has_cert_alerts": &hcldec.AttrSpec{
				Name:     "has_cert_alerts",
				Type:     cty.Bool,
				Required: false,
			},
			"has_kev": &hcldec.AttrSpec{
				Name:     "has_kev",
				Type:     cty.Bool,
				Required: false,
			},
			"has_cert_notes": &hcldec.AttrSpec{
				Name:     "has_cert_notes",
				Type:     cty.Bool,
				Required: false,
			},
			"is_vulnerable": &hcldec.AttrSpec{
				Name:     "is_vulnerable",
				Type:     cty.Bool,
				Required: false,
			},
			"keyword_exact_match": &hcldec.AttrSpec{
				Name:     "keyword_exact_match",
				Type:     cty.Bool,
				Required: false,
			},
			"no_rejected": &hcldec.AttrSpec{
				Name:     "no_rejected",
				Type:     cty.Bool,
				Required: false,
			},
			"limit": &hcldec.AttrSpec{
				Name:     "limit",
				Type:     cty.Number,
				Required: false,
			},
		},
	}
}

func fetchNistNvdCvesData(loader ClientLoadFn) plugin.RetrieveDataFunc {
	return func(ctx context.Context, params *plugin.RetrieveDataParams) (plugin.Data, hcl.Diagnostics) {
		cli, err := parseConfig(params.Config, loader)
		if err != nil {
			return nil, hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Failed to parse configuration",
			}}
		}
		req, err := parseListCVESRequest(params.Args)
		if err != nil {
			return nil, hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Failed to parse arguments",
			}}
		}
		limit := defaultLimit
		if attr := params.Args.GetAttr("limit"); !attr.IsNull() {
			num, _ := attr.AsBigFloat().Int64()
			limit = int(num)
			if limit < minLimit {
				limit = minLimit
			} else if limit > maxLimit {
				limit = maxLimit
			}
		}
		req.ResultsPerPage = limit
		req.StartIndex = 0
		var vulnerabilities plugin.ListData
		for {
			res, err := cli.ListCVES(ctx, req)
			if err != nil {
				return nil, hcl.Diagnostics{{
					Severity: hcl.DiagError,
					Summary:  "Failed to fetch data",
				}}
			}
			for _, v := range res.Vulnerabilities {
				data, err := plugin.ParseDataAny(v)
				if err != nil {
					return nil, hcl.Diagnostics{{
						Severity: hcl.DiagError,
						Summary:  "Failed to parse data",
					}}
				}
				vulnerabilities = append(vulnerabilities, data)
			}
			if res.StartIndex+res.ResultsPerPage >= res.TotalResults {
				break
			}
			req.StartIndex = res.StartIndex + res.ResultsPerPage
		}
		return vulnerabilities, nil
	}
}

func parseConfig(cfg cty.Value, loader ClientLoadFn) (client.Client, error) {
	if cfg.IsNull() {
		return nil, fmt.Errorf("configuration is required")
	}
	apiKey := cfg.GetAttr("api_key")
	if apiKey.IsNull() || apiKey.AsString() == "" {
		return loader(nil), nil
	}
	return loader(client.String(apiKey.AsString())), nil
}

func parseListCVESRequest(args cty.Value) (*client.ListCVESReq, error) {
	if args.IsNull() {
		return nil, fmt.Errorf("arguments are required")
	}
	req := &client.ListCVESReq{}
	if attr := args.GetAttr("last_mod_start_date"); !attr.IsNull() {
		req.LastModStartDate = client.String(attr.AsString())
	}
	if attr := args.GetAttr("last_mod_end_date"); !attr.IsNull() {
		req.LastModEndDate = client.String(attr.AsString())
	}
	if attr := args.GetAttr("pub_start_date"); !attr.IsNull() {
		req.PubStartDate = client.String(attr.AsString())
	}
	if attr := args.GetAttr("pub_end_date"); !attr.IsNull() {
		req.PubEndDate = client.String(attr.AsString())
	}
	if attr := args.GetAttr("cpe_name"); !attr.IsNull() {
		req.CPEName = client.String(attr.AsString())
	}
	if attr := args.GetAttr("cve_id"); !attr.IsNull() {
		req.CVEID = client.String(attr.AsString())
	}
	if attr := args.GetAttr("cvss_v3_metrics"); !attr.IsNull() {
		req.CVSSV3Metrics = client.String(attr.AsString())
	}
	if attr := args.GetAttr("cvss_v3_severity"); !attr.IsNull() {
		req.CVSSV3Severity = client.String(attr.AsString())
	}
	if attr := args.GetAttr("cwe_id"); !attr.IsNull() {
		req.CWEID = client.String(attr.AsString())
	}
	if attr := args.GetAttr("keyword_search"); !attr.IsNull() {
		req.KeywordSearch = client.String(attr.AsString())
	}
	if attr := args.GetAttr("virtual_match_string"); !attr.IsNull() {
		req.VirtualMatchString = client.String(attr.AsString())
	}
	if attr := args.GetAttr("source_identifier"); !attr.IsNull() {
		req.SourceIdentifier = client.String(attr.AsString())
	}
	if attr := args.GetAttr("has_cert_alerts"); !attr.IsNull() {
		req.HasCertAlerts = client.Bool(attr.True())
	}
	if attr := args.GetAttr("has_kev"); !attr.IsNull() {
		req.HasKev = client.Bool(attr.True())
	}
	if attr := args.GetAttr("has_cert_notes"); !attr.IsNull() {
		req.HasCertNotes = client.Bool(attr.True())
	}
	if attr := args.GetAttr("is_vulnerable"); !attr.IsNull() {
		req.IsVulnerable = client.Bool(attr.True())
	}
	if attr := args.GetAttr("keyword_exact_match"); !attr.IsNull() {
		req.KeywordExactMatch = client.Bool(attr.True())
	}
	if attr := args.GetAttr("no_rejected"); !attr.IsNull() {
		req.NoRejected = client.Bool(attr.True())
	}
	return req, nil
}
