package elasticsearch

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"

	es "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/plugin"
)

const (
	defaultBaseURL  = "http://localhost:9200"
	defaultUsername = "elastic"
)

func Plugin(version string) *plugin.Schema {
	return &plugin.Schema{
		Name:    "blackstork/elasticsearch",
		Version: version,
		DataSources: plugin.DataSources{
			"elasticsearch": makeElasticSearchDataSource(),
		},
	}
}

func makeClient(pcfg cty.Value) (*es.Client, error) {
	cfg := &es.Config{
		Addresses: []string{defaultBaseURL},
		Username:  defaultUsername,
	}
	if baseURL := pcfg.GetAttr("base_url"); !baseURL.IsNull() {
		cfg.Addresses = []string{baseURL.AsString()}
	}
	if cloudID := pcfg.GetAttr("cloud_id"); !cloudID.IsNull() {
		cfg.CloudID = cloudID.AsString()
	}
	if apiKeyStr := pcfg.GetAttr("api_key_str"); !apiKeyStr.IsNull() {
		cfg.APIKey = apiKeyStr.AsString()
	}
	if apiKey := pcfg.GetAttr("api_key"); !apiKey.IsNull() {
		list := apiKey.AsValueSlice()
		if len(list) != 2 {
			return nil, fmt.Errorf("api_key must be a list of 2 strings")
		}
		cfg.APIKey = base64.StdEncoding.EncodeToString([]byte(
			fmt.Sprintf("%s:%s", list[0].AsString(), list[1].AsString())),
		)
	}
	if basicAuthUsername := pcfg.GetAttr("basic_auth_username"); !basicAuthUsername.IsNull() {
		cfg.Username = basicAuthUsername.AsString()
	}
	if basicAuthPassword := pcfg.GetAttr("basic_auth_password"); !basicAuthPassword.IsNull() {
		cfg.Password = basicAuthPassword.AsString()
	}
	if bearerAuth := pcfg.GetAttr("bearer_auth"); !bearerAuth.IsNull() {
		cfg.ServiceToken = bearerAuth.AsString()
	}
	if caCerts := pcfg.GetAttr("ca_certs"); !caCerts.IsNull() {
		cfg.CACert = []byte(caCerts.AsString())
	}
	return es.NewClient(*cfg)
}

func getByID(fn esapi.Get, args cty.Value) (plugin.Data, error) {
	index := args.GetAttr("index")
	if index.IsNull() {
		return nil, fmt.Errorf("index is required")
	}
	id := args.GetAttr("id")
	if id.IsNull() {
		return nil, fmt.Errorf("id is required when id is specified")
	}
	opts := []func(*esapi.GetRequest){}
	if fields := args.GetAttr("fields"); !fields.IsNull() {
		fieldSlice := fields.AsValueSlice()
		fieldStrings := make([]string, len(fieldSlice))
		for i, v := range fieldSlice {
			fieldStrings[i] = v.AsString()
		}
		opts = append(opts, fn.WithSource(fieldStrings...))
	}
	res, err := fn(index.AsString(), id.AsString(), opts...)
	if err != nil {
		return nil, err
	}
	if res.IsError() {
		return nil, fmt.Errorf("failed to get document: %s", res.String())
	}
	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read search result: %s", err)
	}
	data, err := plugin.UnmarshalJSONData(raw)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal search result: %s", err)
	}
	return data, nil
}

func search(fn esapi.Search, args cty.Value) (plugin.Data, error) {
	index := args.GetAttr("index")
	if index.IsNull() {
		return nil, fmt.Errorf("index is required")
	}
	opts := []func(*esapi.SearchRequest){
		fn.WithIndex(index.AsString()),
	}
	if queryString := args.GetAttr("query_string"); !queryString.IsNull() {
		opts = append(opts, fn.WithQuery(queryString.AsString()))
	}
	if query := args.GetAttr("query"); !query.IsNull() {
		queryRaw, err := json.Marshal(map[string]any{
			"query": query.AsValueMap(),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to marshal query: %s", err)
		}
		opts = append(opts, fn.WithBody(bytes.NewReader(queryRaw)))
	}
	if fields := args.GetAttr("fields"); !fields.IsNull() {
		fieldSlice := fields.AsValueSlice()
		fieldStrings := make([]string, len(fieldSlice))
		for i, v := range fieldSlice {
			fieldStrings[i] = v.AsString()
		}
		opts = append(opts, fn.WithSource(fieldStrings...))
	}

	res, err := fn(opts...)
	if err != nil {
		return nil, err
	} else if res.IsError() {
		return nil, fmt.Errorf("failed to search: %s", res.String())
	}
	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read search result: %s", err)
	}
	data, err := plugin.UnmarshalJSONData(raw)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal search result: %s", err)
	}
	return data, nil
}
