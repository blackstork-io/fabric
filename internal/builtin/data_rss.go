package builtin

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"sync"
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

	readability "github.com/go-shiori/go-readability"
	"github.com/microcosm-cc/bluemonday"
)

const (
	defaultRequestTimeout           = 30 * time.Second
	defaultUserAgent                = "blackstork-rss/0.0.1"
	defaultOnlyItemsAfterTimeFormat = "2006-01-02T15:04:05Z"
)

// https://techblog.willshouse.com/2012/01/03/most-common-user-agents/
var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:133.0) Gecko/20100101 Firefox/133.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36 Edg/131.0.0.0",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64; rv:133.0) Gecko/20100101 Firefox/133.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:133.0) Gecko/20100101 Firefox/133.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.1.1 Safari/605.1.15",
	"Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:133.0) Gecko/20100101 Firefox/133.0",
}

func getRandUserAgent() string {
	return userAgents[rand.Intn(len(userAgents))]
}

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
				{
					Name:        "fill_in_content",
					Type:        cty.Bool,
					DefaultVal:  cty.BoolVal(false),
					Constraints: constraint.NonNull,
					Doc: `
						If the full content should be added when it's not present in the feed items.
					`,
				},
				{
					Name:        "use_browser_user_agent",
					Type:        cty.Bool,
					DefaultVal:  cty.BoolVal(false),
					Constraints: constraint.NonNull,
					Doc: fmt.Sprintf(`
						If the data source should pretend to be a browser while fetching the feed and the feed items.
						If set to "false", the default user-agent value "%s" will be used.
					`, defaultUserAgent),
				},
				{
					Name:         "fill_in_max_items",
					Type:         cty.Number,
					ExampleVal:   cty.BoolVal(false),
					Constraints:  constraint.NonNull,
					MinInclusive: cty.NumberIntVal(0),
					DefaultVal:   cty.NumberIntVal(10),
					Doc: `
						Maximum number of items to fill the content in per feed.
					`,
				},
				{
					Name:       "only_items_after_time",
					Type:       cty.String,
					ExampleVal: cty.StringVal("2024-12-23T00:00:00Z"),
					Doc: `
						Return only items after a specified date time, in the format "%Y-%m-%dT%H:%M:%S%Z".
					`,
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
		Fetches RSS / Atom / JSON feed from a provided URL.

		The full content of the items can be fetched and added to the feed. The data source supports basic authentication.
		`,
	}
}

func filterItems(ctx context.Context, feed *gofeed.Feed, from time.Time) *gofeed.Feed {
	filteredItems := make([]*gofeed.Item, 0)

	for i := range feed.Items {

		item := feed.Items[i]

		var itemTime *time.Time
		if item.UpdatedParsed != nil {
			itemTime = item.UpdatedParsed
		} else if item.PublishedParsed != nil {
			itemTime = item.PublishedParsed
		}
		if itemTime == nil {
			continue
		} else if itemTime.Before(from) {
			continue
		}

		filteredItems = append(filteredItems, item)
	}
	feed.Items = filteredItems
	return feed
}

func fetchFeedItems(ctx context.Context, feed *gofeed.Feed, userAgent string, itemsCap int) *gofeed.Feed {
	log := slog.Default()
	log = log.With("feed_url", feed.Link, "items_cap", itemsCap)
	log.InfoContext(ctx, "Fetching content for the items in the feed")

	policy := bluemonday.UGCPolicy()

	wg := sync.WaitGroup{}
	count := 0
	for i := range feed.Items {

		if count >= itemsCap {
			log.InfoContext(
				ctx,
				"Max number of items to populate reached for the feed",
				"feed_items_count", len(feed.Items),
			)
			break
		}

		item := feed.Items[i]

		if item.Content != "" {
			log.DebugContext(ctx, "The item already has content, skipping", "item_title", item.Title)
			continue
		}

		wg.Add(1)
		count += 1
		go func(item *gofeed.Item) {
			defer wg.Done()

			_log := log.With("item_title", item.Title, "item_link", item.Link)
			_log.DebugContext(ctx, "Fetching content for the item")

			client := &http.Client{}

			req, err := http.NewRequest("GET", item.Link, nil)
			if err != nil {
				_log.ErrorContext(ctx, "Error while creating a HTTP request for a feed item link", "err", err)
				return
			}
			req.Header.Set("User-Agent", userAgent)

			resp, err := client.Do(req)
			if err != nil {
				_log.ErrorContext(ctx, "Error while fetching a feed item link", "err", err)
				return
			}
			defer resp.Body.Close()

			parsedLink, err := url.Parse(item.Link)
			if err != nil {
				_log.ErrorContext(ctx, "Can't parse the item link", "err", err)
				return
			}

			article, err := readability.FromReader(resp.Body, parsedLink)
			if err != nil {
				_log.ErrorContext(ctx, "Failed to parse a page for a feed item link", "err", err)
				return
			}
			content := policy.Sanitize(article.Content)
			content = strings.TrimSpace(content)
			item.Content = content
		}(item)
	}
	wg.Wait()
	log.InfoContext(ctx, "Feed items has been fetched", "items_fetched_count", count)
	return feed
}

func fetchRSSData(ctx context.Context, params *plugin.RetrieveDataParams) (plugindata.Data, diagnostics.Diag) {
	log := slog.Default()

	fp := gofeed.NewParser()

	url := params.Args.GetAttrVal("url").AsString()

	fillInContent := params.Args.GetAttrVal("fill_in_content").True()
	useBrowserUserAgent := params.Args.GetAttrVal("use_browser_user_agent").True()
	fillInMaxItems, _ := params.Args.GetAttrVal("fill_in_max_items").AsBigFloat().Int64()
	onlyItemsAfterTimeAttr := params.Args.GetAttrVal("only_items_after_time")

	userAgent := defaultUserAgent
	if useBrowserUserAgent {
		userAgent = getRandUserAgent()
	}
	fp.UserAgent = userAgent

	basicAuth := params.Args.Blocks.GetFirstMatching("basic_auth")
	if basicAuth != nil {
		fp.AuthConfig = &gofeed.Auth{
			Username: basicAuth.GetAttrVal("username").AsString(),
			Password: basicAuth.GetAttrVal("password").AsString(),
		}
	}

	ctx, cancel := context.WithTimeout(ctx, defaultRequestTimeout)
	defer cancel()

	log.InfoContext(ctx, "Downloading the feed", "feed_url", url)

	feed, err := fp.ParseURLWithContext(url, ctx)
	if err != nil {
		return nil, diagnostics.Diag{&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("Failed to fetch the feed `%s`", url),
			Detail:   err.Error(),
		}}
	}

	var fromTime time.Time
	if !onlyItemsAfterTimeAttr.IsNull() {
		fromTime, err = time.Parse(defaultOnlyItemsAfterTimeFormat, onlyItemsAfterTimeAttr.AsString())
		if err != nil {
			errorMsg := "Can't parse the value in `only_items_after_time` argument"
			log.ErrorContext(ctx, errorMsg, "err", err)
			return nil, diagnostics.Diag{&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  errorMsg,
				Detail:   err.Error(),
			}}
		}
	}

	if !fromTime.IsZero() {
		oldItemsCount := len(feed.Items)
		feed = filterItems(ctx, feed, fromTime)
		log.InfoContext(
			ctx,
			"Feed items filtered",
			"old_items_count", oldItemsCount,
			"new_items_count", len(feed.Items),
			"only_items_after_time", fromTime,
		)
	}

	if fillInContent {
		feed = fetchFeedItems(ctx, feed, userAgent, int(fillInMaxItems))
		log.InfoContext(ctx, "The content for the feed items downloaded")
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
				"content":     plugindata.String(item.Content),
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
