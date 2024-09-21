package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"github.com/google/go-querystring/query"
)

const (
	authURL    = "https://login.microsoftonline.com"
	defaultURL = "https://management.azure.com"
	version    = "2023-11-01"
)

var scopes = []string{"https://graph.microsoft.com/.default"}

func String(s string) *string {
	return &s
}

func Int(i int) *int {
	return &i
}

type ListIncidentsReq struct {
	SubscriptionID    string  `url:"-"`
	ResourceGroupName string  `url:"-"`
	WorkspaceName     string  `url:"-"`
	Filter            *string `url:"$filter,omitempty"`
	OrderBy           *string `url:"$orderby,omitempty"`
	Top               *int    `url:"$top,omitempty"`
}

type ListIncidentsRes struct {
	Value []any `json:"value"`
}

type GetClientCredentialsTokenReq struct {
	TenantID     string `url:"-"`
	ClientID     string `url:"-"`
	ClientSecret string `url:"-"`
}

type GetClientCredentialsTokenRes struct {
	AccessToken string `json:"access_token"`
}

type client struct {
	url   string
	token string
}

type Client interface {
	UseAuth(token string)
	GetClientCredentialsToken(ctx context.Context, req *GetClientCredentialsTokenReq) (*GetClientCredentialsTokenRes, error)
	ListIncidents(ctx context.Context, req *ListIncidentsReq) (*ListIncidentsRes, error)
}

func New() Client {
	return &client{
		url: defaultURL,
	}
}

func (c *client) prepare(r *http.Request) {
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	q := r.URL.Query()
	q.Add("api-version", version)
	r.URL.RawQuery = q.Encode()
}

func (c *client) UseAuth(token string) {
	c.token = token
}

func (c *client) GetClientCredentialsToken(ctx context.Context, req *GetClientCredentialsTokenReq) (*GetClientCredentialsTokenRes, error) {
	format := "/%s/oauth2/token"
	u, err := url.Parse(authURL + fmt.Sprintf(format, req.TenantID))
	if err != nil {
		return nil, err
	}
	payload := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {req.ClientID},
		"client_secret": {req.ClientSecret},
		"resource":      {defaultURL},
	}
	r, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), strings.NewReader(payload.Encode()))
	if err != nil {
		return nil, err
	}
	c.prepare(r)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := http.Client{
		Timeout: 15 * time.Second,
	}
	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("microsoft sentinels client returned status code: %d", res.StatusCode)
	}
	defer res.Body.Close()
	var data GetClientCredentialsTokenRes
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, err
	}
	return &data, nil
}

func AcquireToken(ctx context.Context, tenantId string, clientId string, cred confidential.Credential) (accessToken string, err error) {
	confidentialClient, err := confidential.New(authURL+"/"+tenantId, clientId, cred)
	if err != nil {
		return
	}
	result, err := confidentialClient.AcquireTokenSilent(ctx, scopes)
	if err != nil {
		// cache miss, authenticate with another AcquireToken... method
		result, err = confidentialClient.AcquireTokenByCredential(ctx, scopes)
		if err != nil {
			return
		}
	}
	accessToken = result.AccessToken
	return
}

func (c *client) ListIncidents(ctx context.Context, req *ListIncidentsReq) (*ListIncidentsRes, error) {
	format := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.OperationalInsights/workspaces/%s/providers/Microsoft.SecurityInsights/incidents"
	u, err := url.Parse(c.url + fmt.Sprintf(format, req.SubscriptionID, req.ResourceGroupName, req.WorkspaceName))
	if err != nil {
		return nil, err
	}
	q, err := query.Values(req)
	if err != nil {
		return nil, err
	}
	u.RawQuery = q.Encode()
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	c.prepare(r)
	client := http.Client{
		Timeout: 15 * time.Second,
	}
	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("microsoft sentinels client returned status code: %d", res.StatusCode)
	}
	defer res.Body.Close()
	var data ListIncidentsRes
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, err
	}
	return &data, nil
}
