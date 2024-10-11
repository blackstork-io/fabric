package client

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"github.com/blackstork-io/fabric/plugin/plugindata"
)

const (
	graphUrl        = "https://graph.microsoft.com"
	defaultPageSize = 20
)

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

func (cli *graphClient) fetchURL(ctx context.Context, requestUrl *url.URL) (result plugindata.Data, err error) {
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, requestUrl.String(), nil)
	if err != nil {
		return
	}
	cli.prepare(r)
	slog.DebugContext(ctx, "Fetching an URL from API", "url", requestUrl.String())
	res, err := cli.client.Do(r)
	if err != nil {
		return
	}
	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read the results: %s", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		slog.ErrorContext(ctx, "Error received from Microsoft Graph API", "status_code", res.StatusCode, "body", string(raw))
		err = fmt.Errorf("Microsoft Graph client returned status code: %d", res.StatusCode)
		return
	}
	result, err = plugindata.UnmarshalJSON(raw)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal results: %s", err)
	}
	return
}

func (cli *graphClient) QueryGraph(
	ctx context.Context,
	endpoint string,
	queryParams url.Values,
	size int,
	onlyobjects bool,
) (result plugindata.Data, err error) {
	objects := make(plugindata.List, 0)

	urlStr := graphUrl + fmt.Sprintf("/%s%s", cli.apiVersion, endpoint)
	requestUrl, err := url.Parse(urlStr)
	if err != nil {
		return
	}

	if queryParams == nil {
		queryParams = url.Values{}
	}

	queryParams.Set("$count", strconv.FormatBool(true))
	queryParams.Set("$skip", strconv.Itoa(0))

	limit := min(size, defaultPageSize)
	queryParams.Set("$top", strconv.Itoa(limit))

	var totalCount int = -1
	var response plugindata.Data

	for {

		requestUrl.RawQuery = queryParams.Encode()

		response, err = cli.fetchURL(ctx, requestUrl)
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

		slog.InfoContext(
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
		queryParams.Set("$skip", strconv.Itoa(len(objects)))

		// Number of objects to get, >= 0
		objectsToGet := max(size-len(objects), 0)
		pageSize := min(defaultPageSize, objectsToGet)
		queryParams.Set("$top", strconv.Itoa(pageSize))
	}

	objectsToReturn := objects[:min(len(objects), size)]

	if onlyobjects {
		return objectsToReturn, nil
	} else {
		data := plugindata.Map{
			"objects":     objectsToReturn,
			"total_count": plugindata.Number(totalCount),
		}
		return data, nil
	}
}

func (cli *graphClient) QueryGraphObject(
	ctx context.Context,
	endpoint string,
) (result plugindata.Data, err error) {
	urlStr := graphUrl + fmt.Sprintf("/%s%s", cli.apiVersion, endpoint)
	requestUrl, err := url.Parse(urlStr)
	if err != nil {
		return
	}
	response, err := cli.fetchURL(ctx, requestUrl)
	if err != nil {
		slog.ErrorContext(ctx, "Error while fetching an object", "url", requestUrl.String(), "error", err)
		return nil, err
	}

	return response, nil
}
