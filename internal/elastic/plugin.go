package elastic

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"

	es "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/internal/elastic/kbclient"
	"github.com/blackstork-io/fabric/plugin"
)

const (
	defaultBaseURL  = "http://localhost:9200"
	defaultUsername = "elastic"
)

type KibanaClientLoaderFn func(url string, apiKey *string) kbclient.Client

var DefaultKibanaClientLoader KibanaClientLoaderFn = kbclient.New

func Plugin(version string, loader KibanaClientLoaderFn) *plugin.Schema {
	if loader == nil {
		loader = DefaultKibanaClientLoader
	}
	return &plugin.Schema{
		Name:    "blackstork/elastic",
		Version: version,
		DataSources: plugin.DataSources{
			"elasticsearch":          makeElasticSearchDataSource(),
			"elastic_security_cases": makeElasticSecurityCasesDataSource(loader),
		},
	}
}

func makeSearchClient(pcfg cty.Value) (*es.Client, error) {
	cfg := &es.Config{
		Username: defaultUsername,
	}
	if pcfg.IsNull() {
		return nil, fmt.Errorf("configuration is required")
	}
	if baseURL := pcfg.GetAttr("base_url"); !baseURL.IsNull() && baseURL.AsString() != "" {
		cfg.Addresses = []string{baseURL.AsString()}
	}
	if cloudID := pcfg.GetAttr("cloud_id"); !cloudID.IsNull() {
		cfg.CloudID = cloudID.AsString()
	}
	if len(cfg.Addresses) == 0 && cfg.CloudID == "" {
		return nil, fmt.Errorf("either one of base_url or cloud_id is required")
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
	if onlyHits := args.GetAttr("only_hits"); onlyHits.IsNull() || onlyHits.True() {
		return extractHits(data)
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
	body := map[string]any{}
	if query := args.GetAttr("query"); !query.IsNull() {
		body["query"] = query.AsValueMap()
	}
	if aggs := args.GetAttr("aggs"); !aggs.IsNull() {
		body["aggs"] = aggs.AsValueMap()
	}
	if len(body) > 0 {
		rawBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal query: %s", err)
		}
		opts = append(opts, fn.WithBody(bytes.NewReader(rawBody)))
	}
	if size := args.GetAttr("size"); !size.IsNull() {
		n, _ := size.AsBigFloat().Int64()
		if n <= 0 {
			return nil, fmt.Errorf("size must be greater than 0")
		}
		opts = append(opts, fn.WithSize(int(n)))
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
	if onlyHits := args.GetAttr("only_hits"); onlyHits.IsNull() || onlyHits.True() {
		return extractHits(data)
	}
	return data, nil
}

func extractHits(data plugin.Data) (plugin.Data, error) {
	m, ok := data.(plugin.MapData)
	if !ok {
		return nil, fmt.Errorf("unexpected search result type: %T", data)
	}
	data, ok = m["hits"]
	if !ok {
		return nil, fmt.Errorf("unexpected search result type: %T", data)
	}
	m, ok = data.(plugin.MapData)
	if !ok {
		return nil, fmt.Errorf("unexpected search result type: %T", data)
	}
	return m["hits"], nil
}
