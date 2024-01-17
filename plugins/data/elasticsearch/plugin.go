package elasticsearch

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/blackstork-io/fabric/plugininterface/v1"
	es "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"
)

const (
	defaultBaseURL  = "http://localhost:9200"
	defaultUsername = "elastic"
)

var Version = semver.MustParse("0.1.0")

type Plugin struct{}

func (Plugin) GetPlugins() []plugininterface.Plugin {
	return []plugininterface.Plugin{
		{
			Namespace: "blackstork",
			Kind:      "data",
			Name:      "elasticsearch",
			Version:   plugininterface.Version(*Version),
			ConfigSpec: &hcldec.ObjectSpec{
				"base_url": &hcldec.AttrSpec{
					Name:     "base_url",
					Type:     cty.String,
					Required: true,
				},
				"cloud_id": &hcldec.AttrSpec{
					Name:     "cloud_id",
					Type:     cty.String,
					Required: false,
				},
				"api_key_str": &hcldec.AttrSpec{
					Name:     "api_key_str",
					Type:     cty.String,
					Required: false,
				},
				"api_key": &hcldec.AttrSpec{
					Name:     "api_key",
					Type:     cty.List(cty.String),
					Required: false,
				},
				"basic_auth_username": &hcldec.AttrSpec{
					Name:     "basic_auth_username",
					Type:     cty.String,
					Required: false,
				},
				"basic_auth_password": &hcldec.AttrSpec{
					Name:     "basic_auth_password",
					Type:     cty.String,
					Required: false,
				},
				"bearer_auth": &hcldec.AttrSpec{
					Name:     "bearer_auth",
					Type:     cty.String,
					Required: false,
				},
				"ca_certs": &hcldec.AttrSpec{
					Name:     "ca_certs",
					Type:     cty.String,
					Required: false,
				},
			},
			InvocationSpec: &hcldec.ObjectSpec{
				"index": &hcldec.AttrSpec{
					Name:     "index",
					Type:     cty.String,
					Required: true,
				},
				"id": &hcldec.AttrSpec{
					Name:     "index",
					Type:     cty.String,
					Required: false,
				},
				"query_string": &hcldec.AttrSpec{
					Name:     "query_string",
					Type:     cty.String,
					Required: false,
				},
				"query": &hcldec.AttrSpec{
					Name:     "query",
					Type:     cty.Map(cty.DynamicPseudoType),
					Required: false,
				},
				"fields": &hcldec.AttrSpec{
					Name:     "fields",
					Type:     cty.List(cty.String),
					Required: false,
				},
			},
		},
	}
}

func (p Plugin) makeClient(pcfg cty.Value) (*es.Client, error) {
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

func (p Plugin) Call(args plugininterface.Args) plugininterface.Result {

	client, err := p.makeClient(args.Config)
	if err != nil {
		return plugininterface.Result{
			Diags: hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Failed to create elasticsearch client",
				Detail:   err.Error(),
			}},
		}
	}
	id := args.Args.GetAttr("id")
	var data map[string]any
	if !id.IsNull() {
		data, err = p.getByID(client.Get, args.Args)
	} else {
		data, err = p.search(client.Search, args.Args)
	}
	if err != nil {
		return plugininterface.Result{
			Diags: hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Failed to get data",
				Detail:   err.Error(),
			}},
		}
	}
	return plugininterface.Result{
		Result: data,
	}
}

func (Plugin) getByID(fn esapi.Get, args cty.Value) (map[string]any, error) {
	index := args.GetAttr("index")
	if index.IsNull() {
		return nil, errors.New("index is required")
	}
	id := args.GetAttr("id")
	if id.IsNull() {
		return nil, errors.New("id is required when id is specified")
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
	var data map[string]any
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal document: %s", err)
	}
	return data, nil
}

func (Plugin) search(fn esapi.Search, args cty.Value) (map[string]any, error) {
	index := args.GetAttr("index")
	if index.IsNull() {
		return nil, errors.New("index is required")
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
	var data map[string]any
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search result: %s", err)
	}
	return data, nil
}
