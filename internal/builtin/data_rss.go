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
)

func makeRSSDataSource() *plugin.DataSource {
	return &plugin.DataSource{
		Doc:      `Fetches an rss or atom feed`,
		Tags:     []string{"rss", "http"},
		DataFunc: fetchRSSData,
		Args: dataspec.ObjectSpec{
			&dataspec.AttrSpec{
				Name:        "url",
				Type:        cty.String,
				ExampleVal:  cty.StringVal("https://www.elastic.co/security-labs/rss/feed.xml"),
				Constraints: constraint.RequiredNonNull,
			},
			&dataspec.BlockSpec{
				Name: "basic_auth",
				Doc: `
					Basic authentication credentials to be used in a HTTP request fetching RSS feed.
				`,
				Nested: dataspec.ObjectSpec{
					&dataspec.AttrSpec{
						Name:        "username",
						Type:        cty.String,
						ExampleVal:  cty.StringVal("user@example.com"),
						Constraints: constraint.RequiredNonNull,
					},
					&dataspec.AttrSpec{
						Name:       "password",
						Type:       cty.String,
						ExampleVal: cty.StringVal("passwd"),
						Doc: `
							Note: you can use function like "from_env_var()" to avoid storing credentials in plaintext
						`,
						Constraints: constraint.RequiredNonNull,
					},
				},
			},
		},
	}
}

func fetchRSSData(ctx context.Context, params *plugin.RetrieveDataParams) (plugin.Data, diagnostics.Diag) {
	fp := gofeed.NewParser()
	url := params.Args.GetAttr("url").AsString()

	basicAuth := params.Args.GetAttr("basic_auth")
	if !basicAuth.IsNull() {
		fp.AuthConfig = &gofeed.Auth{
			Username: basicAuth.GetAttr("username").AsString(),
			Password: basicAuth.GetAttr("password").AsString(),
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
	data := plugin.MapData{
		"title":       plugin.StringData(feed.Title),
		"description": plugin.StringData(feed.Description),
		"link":        plugin.StringData(feed.Link),
		"pub_date":    plugin.StringData(feed.Published),
		"items": plugin.ListData(utils.FnMap(feed.Items, func(item *gofeed.Item) plugin.Data {
			data := plugin.MapData{
				"guid":        plugin.StringData(item.GUID),
				"pub_date":    plugin.StringData(item.Published),
				"title":       plugin.StringData(item.Title),
				"description": plugin.StringData(item.Description),
				"link":        plugin.StringData(item.Link),
			}
			if item.PublishedParsed != nil {
				data["pub_timestamp"] = plugin.NumberData(item.PublishedParsed.Unix())
			}
			return data
		})),
	}
	if feed.PublishedParsed != nil {
		data["pub_timestamp"] = plugin.NumberData(feed.PublishedParsed.Unix())
	}

	return data, nil
}
