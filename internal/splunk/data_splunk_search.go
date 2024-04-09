package splunk

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"time"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/internal/splunk/client"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
)

func makeSplunkSearchDataSchema(loader ClientLoadFn) *plugin.DataSource {
	return &plugin.DataSource{
		Config: dataspec.ObjectSpec{
			&dataspec.AttrSpec{
				Name:     "auth_token",
				Type:     cty.String,
				Required: true,
			},
			&dataspec.AttrSpec{
				Name:     "host",
				Type:     cty.String,
				Required: false,
			},
			&dataspec.AttrSpec{
				Name:     "deployment_name",
				Type:     cty.String,
				Required: false,
			},
		},
		Args: dataspec.ObjectSpec{
			&dataspec.AttrSpec{
				Name:     "search_query",
				Type:     cty.String,
				Required: true,
			},
			&dataspec.AttrSpec{
				Name:     "max_count",
				Type:     cty.Number,
				Required: false,
			},
			&dataspec.AttrSpec{
				Name:     "status_buckets",
				Type:     cty.Number,
				Required: false,
			},
			&dataspec.AttrSpec{
				Name:     "rf",
				Type:     cty.List(cty.String),
				Required: false,
			},
			&dataspec.AttrSpec{
				Name:     "earliest_time",
				Type:     cty.String,
				Required: false,
			},
			&dataspec.AttrSpec{
				Name:     "latest_time",
				Type:     cty.String,
				Required: false,
			},
		},
		DataFunc: fetchSplunkSearchData(loader),
	}
}

func fetchSplunkSearchData(loader ClientLoadFn) plugin.RetrieveDataFunc {
	return func(ctx context.Context, params *plugin.RetrieveDataParams) (plugin.Data, hcl.Diagnostics) {
		cli, err := makeClient(loader, params.Config)
		if err != nil {
			return nil, hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Failed to create client",
				Detail:   err.Error(),
			}}
		}

		result, err := search(cli, ctx, params.Args)
		if err != nil {
			return nil, hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Failed to search",
				Detail:   err.Error(),
			}}
		}
		return result, nil
	}
}

func search(cli client.Client, ctx context.Context, args cty.Value) (plugin.Data, error) {
	id, err := randID()
	if err != nil {
		return nil, err
	}
	req := &client.CreateSearchJobReq{
		ID:       id,
		ExecMode: "blocking",
	}
	if attr := args.GetAttr("search_query"); attr.IsNull() || attr.AsString() == "" {
		return nil, fmt.Errorf("search_query is required")
	} else {
		req.Search = attr.AsString()
	}
	if attr := args.GetAttr("max_count"); !attr.IsNull() {
		n, _ := attr.AsBigFloat().Int64()
		req.MaxCount = client.Int(int(n))
	}
	if attr := args.GetAttr("status_buckets"); !attr.IsNull() {
		n, _ := attr.AsBigFloat().Int64()
		req.StatusBuckets = client.Int(int(n))
	}
	if attr := args.GetAttr("rf"); !attr.IsNull() {
		req.RF = make([]string, attr.LengthInt())
		for i, v := range attr.AsValueSlice() {
			req.RF[i] = v.AsString()
		}
	}
	if attr := args.GetAttr("earliest_time"); !attr.IsNull() {
		req.EarliestTime = client.String(attr.AsString())
	}
	if attr := args.GetAttr("latest_time"); !attr.IsNull() {
		req.LatestTime = client.String(attr.AsString())
	}
	res, err := cli.CreateSearchJob(ctx, req)
	if err != nil {
		return nil, err
	}
	if res.Sid != id {
		return nil, fmt.Errorf("unexpected search job id: %s", res.Sid)
	}
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(1 * time.Second):
			res, err := cli.GetSearchJobByID(ctx, &client.GetSearchJobByIDReq{ID: id})
			if err != nil {
				return nil, err
			}
			if res.DispatchState.Failed() {
				return nil, fmt.Errorf("search job failed: %s", res.DispatchState)
			}
			if res.DispatchState.Done() {
				res, err := cli.GetSearchJobResults(ctx, &client.GetSearchJobResultsReq{
					ID:         id,
					OutputMode: "json",
				})
				if err != nil {
					return nil, err
				}
				result, err := plugin.ParseDataAny(res.Results)
				if err != nil {
					return nil, err
				}
				return result, nil
			}
		}
	}
}

func randID() (string, error) {
	var b [16]byte
	_, err := rand.Read(b[:])
	if err != nil {
		return "", err
	}
	rndStr := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b[:])
	return fmt.Sprintf("fabric_%s", rndStr), nil
}
