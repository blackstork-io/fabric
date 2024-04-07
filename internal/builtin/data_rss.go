package builtin

import (
	"context"
	"time"

	"github.com/hashicorp/hcl/v2"
	"github.com/mmcdole/gofeed"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/pkg/utils"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
)

func makeRSSDataSource() *plugin.DataSource {
	return &plugin.DataSource{
		DataFunc: fetchRSSData,
		Args: dataspec.ObjectSpec{
			&dataspec.AttrSpec{
				Name:       "url",
				Type:       cty.String,
				ExampleVal: cty.StringVal("passwd"),
				Required:   true,
			},
		},
		Config: dataspec.ObjectSpec{
			&dataspec.BlockSpec{
				Name:     "basic_auth",
				Required: false,
				Doc: `
					Authentication parameters used while accessing the rss source.
				`,
				Nested: &dataspec.ObjectSpec{
					&dataspec.AttrSpec{
						Name:       "username",
						Type:       cty.String,
						ExampleVal: cty.StringVal("user@example.com"),
						Required:   true,
					},
					&dataspec.AttrSpec{
						Name:       "password",
						Type:       cty.String,
						ExampleVal: cty.StringVal("passwd"),
						Doc: `
							Note: you can use function like "from_env()" to avoid storing credentials in plaintext
						`,
						Required: true,
					},
				},
			},
		},
	}
}

func fetchRSSData(ctx context.Context, params *plugin.RetrieveDataParams) (plugin.Data, hcl.Diagnostics) {
	fp := gofeed.NewParser()
	url := params.Args.GetAttr("url").AsString()

	basicAuth := params.Config.GetAttr("basic_auth")
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
		return nil, hcl.Diagnostics{&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse the feed",
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
