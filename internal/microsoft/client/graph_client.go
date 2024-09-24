package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const graphUrl = "https://graph.microsoft.com"

type graphClient struct {
	accessToken string
	apiVersion  string
	client      *http.Client
}

func NewGraphClient(accessToken string, apiVersion string) *graphClient {
	return &graphClient{
		accessToken: accessToken,
		apiVersion:  apiVersion,
		client:      &http.Client{},
	}
}

func (cli *graphClient) prepare(r *http.Request) {
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cli.accessToken))
}

func (cli *graphClient) QueryGraph(ctx context.Context, endpoint string, queryParams url.Values) (result interface{}, err error) {
	requestUrl, err := url.Parse(graphUrl + fmt.Sprintf("/%s%s", cli.apiVersion, endpoint))
	if err != nil {
		return
	}
	if queryParams != nil {
		requestUrl.RawQuery = queryParams.Encode()
	}
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, requestUrl.String(), nil)
	if err != nil {
		return
	}
	cli.prepare(r)
	res, err := cli.client.Do(r)
	if err != nil {
		return
	}
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("microsoft graph client returned status code: %d", res.StatusCode)
		return
	}
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}
	return
}
