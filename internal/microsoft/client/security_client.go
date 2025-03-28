package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"github.com/blackstork-io/fabric/plugin/plugindata"
)

const (
	baseURLSecurity         = "https://api.securitycenter.microsoft.com/api"
	defaultPageSizeSecurity = 200
)

type securityClient struct {
	accessToken string
	client      *http.Client
}

func NewSecurityClient(accessToken string) *securityClient {
	return &securityClient{
		accessToken: accessToken,
		client:      &http.Client{},
	}
}

func (client *securityClient) prepare(r *http.Request) {
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.accessToken))
}

func (client *securityClient) getURL(ctx context.Context, requestUrl *url.URL) (result plugindata.Data, err error) {
	slog.DebugContext(ctx, "Sending GET request to an API endpoint", "url", requestUrl.String())
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, requestUrl.String(), nil)
	if err != nil {
		return
	}
	client.prepare(r)
	res, err := client.client.Do(r)
	if err != nil {
		return
	}
	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read the results: %s", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		slog.ErrorContext(
			ctx,
			"Error received from Microsoft Graph API",
			"status_code",
			res.StatusCode,
			"body",
			string(raw),
		)
		err = fmt.Errorf("Microsoft Graph client returned status code: %d", res.StatusCode)
		return
	}
	result, err = plugindata.UnmarshalJSON(raw)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal results: %s", err)
	}
	return
}

func (client *securityClient) postURL(
	ctx context.Context,
	requestUrl *url.URL,
	data plugindata.Data,
) (result plugindata.Data, err error) {
	buff := new(bytes.Buffer)
	err = json.NewEncoder(buff).Encode(data)
	if err != nil {
		return
	}
	slog.DebugContext(ctx, "Sending POST request to an API endpoint", "url", requestUrl.String())

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, requestUrl.String(), buff)
	r.Header.Set("Content-Type", "application/json")
	if err != nil {
		return
	}
	client.prepare(r)
	res, err := client.client.Do(r)
	if err != nil {
		return
	}
	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read the results: %s", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		slog.ErrorContext(ctx, "Error received from API", "status_code", res.StatusCode, "body", string(raw))
		if res.StatusCode == http.StatusNotFound {
			return nil, nil
		}
		err = fmt.Errorf("API returned status code %d", res.StatusCode)
		return
	}

	result, err = plugindata.UnmarshalJSON(raw)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal results: %s", err)
	}
	return result, nil
}

func (client *securityClient) QueryObjects(
	ctx context.Context,
	endpoint string,
	queryParams url.Values,
	size int,
) (result plugindata.List, err error) {
	objects := make(plugindata.List, 0)

	urlStr := baseURLSecurity + endpoint
	requestUrl, err := url.Parse(urlStr)
	if err != nil {
		return
	}

	if queryParams == nil {
		queryParams = url.Values{}
	}

	limit := min(size, defaultPageSizeSecurity)
	queryParams.Set("$top", strconv.Itoa(limit))

	requestUrl.RawQuery = queryParams.Encode()

	var totalCount int = -1
	var response plugindata.Data

	for {
		slog.DebugContext(ctx, "Fetching a page from Microsoft Graph API", "url", requestUrl.String())
		response, err = client.getURL(ctx, requestUrl)
		if err != nil {
			slog.ErrorContext(ctx, "Error while fetching objects", "url", requestUrl.String(), "error", err)
			return nil, err
		}

		resultMap, ok := response.(plugindata.Map)
		if !ok {
			return nil, fmt.Errorf("unexpected result type: %T", response)
		}

		countRaw, ok := resultMap["@odata.count"]
		if ok {
			totalCount = int(countRaw.(plugindata.Number))
		}

		objectsPageRaw, ok := resultMap["value"]
		if !ok {
			break
		}

		objectsPage, ok := objectsPageRaw.(plugindata.List)
		if !ok {
			return nil, fmt.Errorf("unexpected value type: %T", objectsPageRaw)
		}

		if len(objectsPage) == 0 {
			break
		}

		slog.DebugContext(
			ctx, "Objects fetched from Microsoft Graph API",
			"fetched_overall", len(objects),
			"fetched", len(objectsPage),
			"total_available", totalCount,
			"to_fetch_overall", size,
		)

		objects = append(objects, objectsPage...)
		if len(objects) >= size {
			break
		}

		nextLink, ok := resultMap["@odata.nextLink"]
		if !ok && nextLink == nil {
			break
		}
		requestUrlRaw, ok := nextLink.(plugindata.String)
		if !ok {
			return nil, fmt.Errorf("unexpected value type for `@odata.nextLink`: %T", requestUrlRaw)
		}
		requestUrl, err = url.Parse(string(requestUrlRaw))
		if err != nil {
			slog.DebugContext(ctx, "Can't parse the next link in Microsoft Graph API response", "value", requestUrlRaw)
			return nil, err
		}
	}

	objectsToReturn := objects[:min(len(objects), size)]
	return objectsToReturn, nil
}

func (client *securityClient) QueryObject(
	ctx context.Context,
	endpoint string,
	queryParams url.Values,
) (result plugindata.Data, err error) {
	urlStr := baseURLSecurity + endpoint
	requestUrl, err := url.Parse(urlStr)
	if err != nil {
		return
	}

	if queryParams == nil {
		queryParams = url.Values{}
	}
	requestUrl.RawQuery = queryParams.Encode()

	response, err := client.getURL(ctx, requestUrl)
	if err != nil {
		slog.ErrorContext(ctx, "Error while fetching an object", "url", requestUrl.String(), "error", err)
		return nil, err
	}

	return response, nil
}

func (client *securityClient) RunAdvancedQuery(ctx context.Context, query string) (result plugindata.Data, err error) {
	urlStr := baseURLSecurity + "/advancedqueries/run"
	requestUrl, err := url.Parse(urlStr)
	if err != nil {
		return
	}

	body := plugindata.Map{
		"Query": plugindata.String(query),
	}

	response, err := client.postURL(ctx, requestUrl, body)
	if err != nil {
		slog.ErrorContext(ctx, "Error while submitting an advanced query", "url", requestUrl.String(), "error", err, "query", query)
		return nil, err
	}

	return response, nil
}
