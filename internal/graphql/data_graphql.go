package graphql

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/plugin"
)

func makeGraphQLDataSource() *plugin.DataSource {
	return &plugin.DataSource{
		Config: hcldec.ObjectSpec{
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
		Args: hcldec.ObjectSpec{
			"query": &hcldec.AttrSpec{
				Name:     "query",
				Type:     cty.String,
				Required: true,
			},
		},
		DataFunc: fetchGraphQLData,
	}
}

func fetchGraphQLData(ctx context.Context, params *plugin.RetrieveDataParams) (plugin.Data, hcl.Diagnostics) {
	url := params.Config.GetAttr("url")
	if url.IsNull() || url.AsString() == "" {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse config",
			Detail:   "url is required",
		}}
	}
	authToken := params.Config.GetAttr("auth_token")
	if authToken.IsNull() {
		authToken = cty.StringVal("")
	}
	query := params.Args.GetAttr("query")
	if query.IsNull() || query.AsString() == "" {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse arguments",
			Detail:   "query is required",
		}}
	}

	result, err := queryGraphQL(ctx, url.AsString(), query.AsString(), authToken.AsString())
	if err != nil {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to execute query",
			Detail:   err.Error(),
		}}
	}

	return result, nil
}

type requestData struct {
	Query string `json:"query"`
}

func queryGraphQL(ctx context.Context, url, query, authToken string) (plugin.Data, error) {
	data, err := json.Marshal(requestData{Query: query})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
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
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	dst, err := plugin.UnmarshalJSONData(raw)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return dst, nil
}
