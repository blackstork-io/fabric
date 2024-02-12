package openai

import (
	"fmt"

	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/internal/openai/client"
	"github.com/blackstork-io/fabric/plugin"
)

const (
	defaultModel   = "gpt-3.5-turbo"
	queryResultKey = "query_result"
)

type ClientLoadFn func(opts ...client.Option) client.Client

var DefaultClientLoader ClientLoadFn = client.New

func Plugin(version string, loader ClientLoadFn) *plugin.Schema {
	if loader == nil {
		loader = DefaultClientLoader
	}
	return &plugin.Schema{
		Name:    "blackstork/openai",
		Version: version,
		ContentProviders: plugin.ContentProviders{
			"openai_text": makeOpenAITextContentSchema(loader),
		},
	}
}

func makeClient(loader ClientLoadFn, cfg cty.Value) (client.Client, error) {
	opts := []client.Option{}
	apiKey := cfg.GetAttr("api_key")
	if apiKey.IsNull() || apiKey.AsString() == "" {
		return nil, fmt.Errorf("api_key is required in configuration")
	}
	opts = append(opts, client.WithAPIKey(apiKey.AsString()))
	orgID := cfg.GetAttr("organization_id")
	if !orgID.IsNull() && orgID.AsString() != "" {
		opts = append(opts, client.WithOrgID(orgID.AsString()))
	}
	return loader(opts...), nil
}
