package builtin

import (
	"context"
	"time"

	"github.com/hashicorp/hcl/v2"
	"github.com/mmcdole/gofeed"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/pkg/utils"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/dataspec/constraint"
	"github.com/blackstork-io/fabric/plugin/plugindata"
)

func makeRSSDataSource() *plugin.DataSource {
	return &plugin.DataSource{
		Tags:     []string{"rss", "http"},
		DataFunc: fetchRSSData,
		Args: &dataspec.RootSpec{
			Attrs: []*dataspec.AttrSpec{
				{
					Name:        "url",
					Type:        cty.String,
					ExampleVal:  cty.StringVal("https://www.elastic.co/security-labs/rss/feed.xml"),
					Constraints: constraint.RequiredNonNull,
				},
			},
			Blocks: []*dataspec.BlockSpec{
				{
					Header: dataspec.HeadersSpec{
						dataspec.ExactMatcher{"basic_auth"},
					},
					Doc: `
						Basic authentication credentials to be used in a HTTP request fetching RSS feed.
					`,
					Attrs: []*dataspec.AttrSpec{
						{
							Name:        "username",
							Type:        cty.String,
							ExampleVal:  cty.StringVal("user@example.com"),
							Constraints: constraint.RequiredNonNull,
						},
						{
							Name:       "password",
							Type:       cty.String,
							ExampleVal: cty.StringVal("passwd"),
							Doc: `
								Note: avoid storing credentials in the templates. Use environment variables instead.
							`,
							Constraints: constraint.RequiredNonNull,
						},
					},
				},
			},
		},
		Doc: `
		Fetches RSS / Atom feed from a URL.

		The data source supports basic authentication.
		`,
	}
}

func fetchRSSData(ctx context.Context, params *plugin.RetrieveDataParams) (plugindata.Data, diagnostics.Diag) {
	fp := gofeed.NewParser()
	url := params.Args.GetAttrVal("url").AsString()

	basicAuth := params.Args.Blocks.GetFirstMatching("basic_auth")
	if basicAuth != nil {
		fp.AuthConfig = &gofeed.Auth{
			Username: basicAuth.GetAttrVal("username").AsString(),
			Password: basicAuth.GetAttrVal("password").AsString(),
		}
	}

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	feed, err := fp.ParseURLWithContext(url, ctx)
	if err != nil {
		return nil, diagnostics.Diag{&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Failed to fetch the feed",
			Detail:   err.Error(),
		}}
	}
	data := plugindata.Map{
		"title":       plugindata.String(feed.Title),
		"description": plugindata.String(feed.Description),
		"link":        plugindata.String(feed.Link),
		"pub_date":    plugindata.String(feed.Published),
		"items": plugindata.List(utils.FnMap(feed.Items, func(item *gofeed.Item) plugindata.Data {
			data := plugindata.Map{
				"guid":        plugindata.String(item.GUID),
				"pub_date":    plugindata.String(item.Published),
				"title":       plugindata.String(item.Title),
				"description": plugindata.String(item.Description),
				"link":        plugindata.String(item.Link),
			}
			if item.PublishedParsed != nil {
				data["pub_timestamp"] = plugindata.Number(item.PublishedParsed.Unix())
			}
			return data
		})),
	}
	if feed.PublishedParsed != nil {
		data["pub_timestamp"] = plugindata.Number(feed.PublishedParsed.Unix())
	}

	return data, nil
}
