package graphql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Masterminds/semver/v3"
	"github.com/blackstork-io/fabric/plugininterface/v1"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"
)

var Version = semver.MustParse("0.1.0")

type Plugin struct{}

func (Plugin) GetPlugins() []plugininterface.Plugin {
	return []plugininterface.Plugin{
		{
			Namespace: "blackstork",
			Kind:      "data",
			Name:      "graphql",
			Version:   plugininterface.Version(*Version),
			ConfigSpec: &hcldec.ObjectSpec{
				"url": &hcldec.AttrSpec{
					Name:     "url",
					Type:     cty.String,
					Required: true,
				},
				"auth_token": &hcldec.AttrSpec{
					Name:     "auth_token",
					Type:     cty.String,
					Required: false,
				},
			},
			InvocationSpec: &hcldec.ObjectSpec{
				"query": &hcldec.AttrSpec{
					Name:     "path",
					Type:     cty.String,
					Required: true,
				},
			},
		},
	}
}

func (Plugin) parseConfig(cfg cty.Value) (string, string, error) {
	url := cfg.GetAttr("url")
	if url.IsNull() || url.AsString() == "" {
		return "", "", fmt.Errorf("url is required")
	}
	authToken := cfg.GetAttr("auth_token")
	if authToken.IsNull() {
		authToken = cty.StringVal("")
	}
	return url.AsString(), authToken.AsString(), nil
}

func (Plugin) parseArgs(args cty.Value) (string, error) {
	query := args.GetAttr("query")
	if query.IsNull() || query.AsString() == "" {
		return "", fmt.Errorf("query is required")
	}
	return query.AsString(), nil
}

func (p Plugin) Call(args plugininterface.Args) plugininterface.Result {
	url, authToken, err := p.parseConfig(args.Config)
	if err != nil {
		return plugininterface.Result{
			Diags: hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Failed to parse config",
				Detail:   err.Error(),
			}},
		}
	}
	query, err := p.parseArgs(args.Args)
	if err != nil {
		return plugininterface.Result{
			Diags: hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Failed to parse arguments",
				Detail:   err.Error(),
			}},
		}
	}

	result, err := p.query(url, query, authToken)
	if err != nil {
		return plugininterface.Result{
			Diags: hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Failed to execute query",
				Detail:   err.Error(),
			}},
		}
	}

	return plugininterface.Result{
		Result: result,
	}
}

type requestData struct {
	Query string `json:"query"`
}

func (Plugin) query(url, query, authToken string) (any, error) {
	data, err := json.Marshal(requestData{Query: query})
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	// Set the appropriate headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if authToken != "" {
		req.Header.Set("Authorization", "Bearer "+authToken)
	}
	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	var result any
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}
	return result, nil
}
