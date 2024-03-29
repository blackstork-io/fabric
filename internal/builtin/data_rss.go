package builtin

import (
	"context"
	"time"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/mmcdole/gofeed"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/pkg/utils"
	"github.com/blackstork-io/fabric/plugin"
)

func makeRSSDataSource() *plugin.DataSource {
	return &plugin.DataSource{
		DataFunc: fetchRSSData,
		Args: hcldec.ObjectSpec{
			"url": &hcldec.AttrSpec{
				Name:     "url",
				Type:     cty.String,
				Required: true,
			},
		},
		Config: hcldec.ObjectSpec{
			"basic_auth": &hcldec.BlockListSpec{
				TypeName: "basic_auth",
				Nested: &hcldec.ObjectSpec{
					"username": &hcldec.AttrSpec{
						Name:     "username",
						Type:     cty.String,
						Required: true,
					},
					"password": &hcldec.AttrSpec{
						Name:     "password",
						Type:     cty.String,
						Required: true,
					},
				},
				MinItems: 0,
				MaxItems: 1,
			},
		},
	}
}

func fetchRSSData(ctx context.Context, params *plugin.RetrieveDataParams) (plugin.Data, hcl.Diagnostics) {
	fp := gofeed.NewParser()
	url := params.Args.GetAttr("url").AsString()

	basicAuthList := params.Config.GetAttr("basic_auth")
	if basicAuthList.LengthInt() == 1 {
		basicAuthObj := basicAuthList.Index(cty.NumberIntVal(0))
		fp.AuthConfig = &gofeed.Auth{
			Username: basicAuthObj.GetAttr("username").AsString(),
			Password: basicAuthObj.GetAttr("password").AsString(),
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
		"items": plugin.ListData(utils.FnMap(func(item *gofeed.Item) plugin.Data {
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
		}, feed.Items)),
	}
	if feed.PublishedParsed != nil {
		data["pub_timestamp"] = plugin.NumberData(feed.PublishedParsed.Unix())
	}

	return data, nil
}
