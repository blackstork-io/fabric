package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ClientTestSuite struct {
	suite.Suite
	ctx    context.Context
	cancel context.CancelFunc
}

func (s *ClientTestSuite) SetupTest() {
	s.ctx, s.cancel = context.WithCancel(context.Background())
}

func (s *ClientTestSuite) TearDownTest() {
	s.cancel()
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}

func (s *ClientTestSuite) mock(fn http.HandlerFunc, token string) (azureClient, *httptest.Server) {
	srv := httptest.NewServer(fn)
	cli := azureClient{
		accessToken: token,
		baseURL:     srv.URL,
		client:      &http.Client{},
	}
	return cli, srv
}

func (s *ClientTestSuite) TestPrepare() {
	client, srv := s.mock(func(w http.ResponseWriter, r *http.Request) {
		s.Equal("Bearer test_token", r.Header.Get("Authorization"))
		s.Equal("2023-11-01", r.URL.Query().Get("api-version"))
	}, "test_token")
	defer srv.Close()
	client.QueryObjects(s.ctx, "/tmp", url.Values{}, 1)
}

// func (s *ClientTestSuite) TestListIncidents() {
// 	client, srv := s.mock(func(w http.ResponseWriter, r *http.Request) {
// 		expectPath := "/subscriptions/test_subscription_id/resourceGroups/test_resource_group/providers" +
// 			"/Microsoft.OperationalInsights/workspaces/test_workspace/providers/Microsoft.SecurityInsights/incidents"
// 		s.Equal(expectPath, r.URL.Path)
// 		s.Equal(http.MethodGet, r.Method)
// 		s.Equal("Bearer test_token", r.Header.Get("Authorization"))
// 		s.Equal("10", r.URL.Query().Get("$top"))
// 		s.Equal("test_filter", r.URL.Query().Get("$filter"))
// 		s.Equal("test_order_by", r.URL.Query().Get("$orderby"))
// 		w.Write([]byte(`{
// 			"value": [
// 				{
// 					"any": "data"
// 				}
// 			]
// 		}`))
// 	}, "test_token")
// 	defer srv.Close()
// 	req := ListIncidentsReq{
// 		SubscriptionID:    "test_subscription_id",
// 		ResourceGroupName: "test_resource_group",
// 		WorkspaceName:     "test_workspace",
// 		Filter:            String("test_filter"),
// 		OrderBy:           String("test_order_by"),
// 		Top:               Int(10),
// 	}
// 	result, err := client.ListIncidents(s.ctx, &req)
// 	s.NoError(err)
// 	s.Equal(&ListIncidentsRes{
// 		Value: []any{
// 			map[string]any{
// 				"any": "data",
// 			},
// 		},
// 	}, result)
// }

func (s *ClientTestSuite) TestGetAllReportsError() {
	client, srv := s.mock(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}, "test_token")
	defer srv.Close()
	_, err := client.QueryObjects(s.ctx, "/tmp", url.Values{}, 1)
	s.Error(err)
}

func (s *ClientTestSuite) TestBaseURL() {
	cli := NewAzureClient("dummy-access-token")
	s.Equal(baseURLAzure, cli.baseURL)
}
