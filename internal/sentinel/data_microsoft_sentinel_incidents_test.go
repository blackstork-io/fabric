package sentinel

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/internal/sentinel/client"
	client_mocks "github.com/blackstork-io/fabric/mocks/internalpkg/sentinel/client"
	"github.com/blackstork-io/fabric/plugin"
)

type SentinelIncidentsDataSourceTestSuite struct {
	suite.Suite
	schema    *plugin.DataSource
	ctx       context.Context
	cli       *client_mocks.Client
	storedTkn string
}

func TestSentinelIncidentsDataSourceTestSuite(t *testing.T) {
	suite.Run(t, new(SentinelIncidentsDataSourceTestSuite))
}

func (s *SentinelIncidentsDataSourceTestSuite) SetupSuite() {
	s.schema = makeMicrosoftSentinelIncidentsDataSource(func() client.Client {
		return s.cli
	})
	s.ctx = context.Background()
}

func (s *SentinelIncidentsDataSourceTestSuite) SetupTest() {
	s.cli = &client_mocks.Client{}
}

func (s *SentinelIncidentsDataSourceTestSuite) TearDownTest() {
	s.cli.AssertExpectations(s.T())
}

func (s *SentinelIncidentsDataSourceTestSuite) TestSchema() {
	s.Require().NotNil(s.schema)
	s.NotNil(s.schema.Config)
	s.NotNil(s.schema.Args)
	s.NotNil(s.schema.DataFunc)
}

func (s *SentinelIncidentsDataSourceTestSuite) Testlimit() {
	s.cli.On("GetClientCredentialsToken", mock.Anything, &client.GetClientCredentialsTokenReq{
		TenantID:     "test_tenant_id",
		ClientID:     "test_client_id",
		ClientSecret: "test_client_secret",
	}).Return(&client.GetClientCredentialsTokenRes{
		AccessToken: "test_token",
	}, nil)
	s.cli.On("UseAuth", "test_token").Run(func(args mock.Arguments) {
		s.storedTkn = args.Get(0).(string)
	}).Return()
	s.cli.On("ListIncidents", mock.Anything, &client.ListIncidentsReq{
		SubscriptionID:    "test_subscription_id",
		ResourceGroupName: "test_resource_group_name",
		WorkspaceName:     "test_workspace_name",
		Top:               client.Int(123),
	}).Return(&client.ListIncidentsRes{
		Value: []any{
			map[string]any{
				"any": "data",
			},
		},
	}, nil)
	res, diags := s.schema.DataFunc(s.ctx, &plugin.RetrieveDataParams{
		Config: cty.ObjectVal(map[string]cty.Value{
			"tenant_id":           cty.StringVal("test_tenant_id"),
			"client_id":           cty.StringVal("test_client_id"),
			"client_secret":       cty.StringVal("test_client_secret"),
			"subscription_id":     cty.StringVal("test_subscription_id"),
			"resource_group_name": cty.StringVal("test_resource_group_name"),
			"workspace_name":      cty.StringVal("test_workspace_name"),
		}),
		Args: cty.ObjectVal(map[string]cty.Value{
			"limit":    cty.NumberIntVal(123),
			"filter":   cty.NullVal(cty.String),
			"order_by": cty.NullVal(cty.String),
		}),
	})
	s.Equal("test_token", s.storedTkn)
	s.Len(diags, 0)
	s.Equal(plugin.ListData{
		plugin.MapData{
			"any": plugin.StringData("data"),
		},
	}, res)
}

func (s *SentinelIncidentsDataSourceTestSuite) TestFull() {
	s.cli.On("GetClientCredentialsToken", mock.Anything, &client.GetClientCredentialsTokenReq{
		TenantID:     "test_tenant_id",
		ClientID:     "test_client_id",
		ClientSecret: "test_client_secret",
	}).Return(&client.GetClientCredentialsTokenRes{
		AccessToken: "test_token",
	}, nil)
	s.cli.On("UseAuth", "test_token").Run(func(args mock.Arguments) {
		s.storedTkn = args.Get(0).(string)
	}).Return()
	s.cli.On("ListIncidents", mock.Anything, &client.ListIncidentsReq{
		SubscriptionID:    "test_subscription_id",
		ResourceGroupName: "test_resource_group_name",
		WorkspaceName:     "test_workspace_name",
		Filter:            client.String("test_filter"),
		OrderBy:           client.String("test_order_by"),
		Top:               client.Int(10),
	}).Return(&client.ListIncidentsRes{
		Value: []any{
			map[string]any{
				"any": "data",
			},
		},
	}, nil)
	res, diags := s.schema.DataFunc(s.ctx, &plugin.RetrieveDataParams{
		Config: cty.ObjectVal(map[string]cty.Value{
			"tenant_id":           cty.StringVal("test_tenant_id"),
			"client_id":           cty.StringVal("test_client_id"),
			"client_secret":       cty.StringVal("test_client_secret"),
			"subscription_id":     cty.StringVal("test_subscription_id"),
			"resource_group_name": cty.StringVal("test_resource_group_name"),
			"workspace_name":      cty.StringVal("test_workspace_name"),
		}),
		Args: cty.ObjectVal(map[string]cty.Value{
			"filter":   cty.StringVal("test_filter"),
			"order_by": cty.StringVal("test_order_by"),
			"limit":    cty.NumberIntVal(10),
		}),
	})
	s.Equal("test_token", s.storedTkn)
	s.Len(diags, 0)
	s.Equal(plugin.ListData{
		plugin.MapData{
			"any": plugin.StringData("data"),
		},
	}, res)
}

func (s *SentinelIncidentsDataSourceTestSuite) TestError() {
	errTest := fmt.Errorf("test_error")
	s.cli.On("GetClientCredentialsToken", mock.Anything, &client.GetClientCredentialsTokenReq{
		TenantID:     "test_tenant_id",
		ClientID:     "test_client_id",
		ClientSecret: "test_client_secret",
	}).Return(&client.GetClientCredentialsTokenRes{
		AccessToken: "test_token",
	}, nil)
	s.cli.On("UseAuth", "test_token").Run(func(args mock.Arguments) {
		s.storedTkn = args.Get(0).(string)
	}).Return()
	s.cli.On("ListIncidents", mock.Anything, &client.ListIncidentsReq{
		SubscriptionID:    "test_subscription_id",
		ResourceGroupName: "test_resource_group_name",
		WorkspaceName:     "test_workspace_name",
		Top:               client.Int(10),
	}).Return(nil, errTest)
	_, diags := s.schema.DataFunc(s.ctx, &plugin.RetrieveDataParams{
		Config: cty.ObjectVal(map[string]cty.Value{
			"tenant_id":           cty.StringVal("test_tenant_id"),
			"client_id":           cty.StringVal("test_client_id"),
			"client_secret":       cty.StringVal("test_client_secret"),
			"subscription_id":     cty.StringVal("test_subscription_id"),
			"resource_group_name": cty.StringVal("test_resource_group_name"),
			"workspace_name":      cty.StringVal("test_workspace_name"),
		}),
		Args: cty.ObjectVal(map[string]cty.Value{
			"limit":    cty.NumberIntVal(10),
			"filter":   cty.NullVal(cty.String),
			"order_by": cty.NullVal(cty.String),
		}),
	})
	s.Equal("test_token", s.storedTkn)
	s.Len(diags, 1)
	s.Equal("Unable to list Microsoft Sentinel incidents", diags[0].Summary)
	s.Equal(errTest.Error(), diags[0].Detail)
}
